package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"saleapp/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id uuid.UUID) (*models.Order, error)
	GetByOrderNumber(orderNumber string) (*models.Order, error)
	Update(order *models.Order) error
	Delete(id uuid.UUID) error
	List(limit, offset int, customerID *uuid.UUID, status string, startDate, endDate *time.Time) ([]models.Order, int64, error)
	GetSalesByDate(startDate, endDate time.Time) ([]models.Order, error)
	GetTopSellingProducts(startDate, endDate time.Time, limit int) ([]struct {
		ProductID uuid.UUID `gorm:"column:product_id"`
		TotalSold int64     `gorm:"column:total_sold"`
	}, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Customer").Preload("User").Preload("Items.Product").
		First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, "order_number = ?", orderNumber).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}

func (r *orderRepository) List(limit, offset int, customerID *uuid.UUID, status string, startDate, endDate *time.Time) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{})

	if customerID != nil {
		query = query.Where("customer_id = ?", *customerID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	query.Count(&total)

	err := query.Preload("Customer").Preload("User").Preload("Items.Product").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) GetSalesByDate(startDate, endDate time.Time) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("status IN ? AND created_at >= ? AND created_at <= ?",
		[]models.OrderStatus{models.StatusCompleted, models.StatusPending},
		startDate, endDate).Preload("Items").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetTopSellingProducts(startDate, endDate time.Time, limit int) ([]struct {
	ProductID uuid.UUID `gorm:"column:product_id"`
	TotalSold int64     `gorm:"column:total_sold"`
}, error) {
	var results []struct {
		ProductID uuid.UUID `gorm:"column:product_id"`
		TotalSold int64     `gorm:"column:total_sold"`
	}

	err := r.db.Model(&models.OrderItem{}).
		Select("order_items.product_id, SUM(order_items.quantity) as total_sold").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status IN ? AND orders.created_at >= ? AND orders.created_at <= ?",
			[]models.OrderStatus{models.StatusCompleted, models.StatusPending},
			startDate, endDate).
		Group("order_items.product_id").
		Order("total_sold DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}
