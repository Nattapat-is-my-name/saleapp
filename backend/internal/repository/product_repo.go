package repository

import (
	"errors"

	"github.com/google/uuid"
	"saleapp/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uuid.UUID) (*models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uuid.UUID) error
	List(limit, offset int, search string, categoryID *uuid.UUID, isActive *bool) ([]models.Product, int64, error)
	UpdateStock(id uuid.UUID, quantity int) error
	GetLowStock(threshold int) ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "sku = ?", sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}

func (r *productRepository) List(limit, offset int, search string, categoryID *uuid.UUID, isActive *bool) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR sku ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	query.Count(&total)

	err := query.Preload("Category").Limit(limit).Offset(offset).Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) UpdateStock(id uuid.UUID, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *productRepository) GetLowStock(threshold int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("stock > 0 AND stock <= ?", threshold).Find(&products).Error
	return products, err
}
