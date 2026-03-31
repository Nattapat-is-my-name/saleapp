package service

import (
	"github.com/google/uuid"
	"saleapp/internal/dto/request"
	"saleapp/internal/models"
	"saleapp/internal/repository"
	"saleapp/pkg/errors"
)

type CustomerService interface {
	Create(req *request.CreateCustomerRequest) (*models.Customer, error)
	GetByID(id uuid.UUID) (*models.Customer, error)
	Update(id uuid.UUID, req *request.UpdateCustomerRequest) (*models.Customer, error)
	Delete(id uuid.UUID) error
	List(req *request.ListCustomersRequest) ([]models.Customer, int64, error)
}

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
	}
}

func (s *customerService) Create(req *request.CreateCustomerRequest) (*models.Customer, error) {
	if req.Email != "" {
		existing, err := s.customerRepo.GetByEmail(req.Email)
		if err != nil {
			return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to check existing customer")
		}
		if existing != nil {
			return nil, errors.Wrap(errors.ErrDuplicateEntry, "EMAIL_EXISTS", "Email already exists")
		}
	}

	customer := &models.Customer{
		Email:     req.Email,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Address:   req.Address,
		Notes:     req.Notes,
	}

	if err := s.customerRepo.Create(customer); err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to create customer")
	}

	return customer, nil
}

func (s *customerService) GetByID(id uuid.UUID) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get customer")
	}
	if customer == nil {
		return nil, errors.ErrNotFound
	}
	return customer, nil
}

func (s *customerService) Update(id uuid.UUID, req *request.UpdateCustomerRequest) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get customer")
	}
	if customer == nil {
		return nil, errors.ErrNotFound
	}

	if req.Email != nil {
		existing, err := s.customerRepo.GetByEmail(*req.Email)
		if err != nil {
			return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to check email")
		}
		if existing != nil && existing.ID != id {
			return nil, errors.Wrap(errors.ErrDuplicateEntry, "EMAIL_EXISTS", "Email already exists")
		}
		customer.Email = *req.Email
	}

	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.FirstName != nil {
		customer.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		customer.LastName = *req.LastName
	}
	if req.Address != nil {
		customer.Address = *req.Address
	}
	if req.Notes != nil {
		customer.Notes = *req.Notes
	}

	if err := s.customerRepo.Update(customer); err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to update customer")
	}

	return customer, nil
}

func (s *customerService) Delete(id uuid.UUID) error {
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return errors.Wrap(err, "INTERNAL_ERROR", "Failed to get customer")
	}
	if customer == nil {
		return errors.ErrNotFound
	}

	if err := s.customerRepo.Delete(id); err != nil {
		return errors.Wrap(err, "INTERNAL_ERROR", "Failed to delete customer")
	}

	return nil
}

func (s *customerService) List(req *request.ListCustomersRequest) ([]models.Customer, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	customers, total, err := s.customerRepo.List(limit, offset, req.Search)
	if err != nil {
		return nil, 0, errors.Wrap(err, "INTERNAL_ERROR", "Failed to list customers")
	}

	return customers, total, nil
}
