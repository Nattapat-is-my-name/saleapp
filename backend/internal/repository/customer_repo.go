package repository

import (
	"errors"

	"github.com/google/uuid"
	"saleapp/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	GetByID(id uuid.UUID) (*models.Customer, error)
	GetByEmail(email string) (*models.Customer, error)
	Update(customer *models.Customer) error
	Delete(id uuid.UUID) error
	List(limit, offset int, search string) ([]models.Customer, int64, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) GetByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetByEmail(email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) List(limit, offset int, search string) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	query := r.db.Model(&models.Customer{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", 
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	query.Count(&total)

	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
