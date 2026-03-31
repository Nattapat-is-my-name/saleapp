package service

import (
	"github.com/google/uuid"
	"saleapp/internal/dto/request"
	"saleapp/internal/models"
	"saleapp/internal/repository"
	"saleapp/pkg/errors"
)

type ProductService interface {
	Create(req *request.CreateProductRequest) (*models.Product, error)
	GetByID(id uuid.UUID) (*models.Product, error)
	Update(id uuid.UUID, req *request.UpdateProductRequest) (*models.Product, error)
	Delete(id uuid.UUID) error
	List(req *request.ListProductsRequest) ([]models.Product, int64, error)
	GetLowStock(threshold int) ([]models.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) Create(req *request.CreateProductRequest) (*models.Product, error) {
	existing, err := s.productRepo.GetBySKU(req.SKU)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to check existing product")
	}
	if existing != nil {
		return nil, errors.Wrap(errors.ErrDuplicateEntry, "SKU_EXISTS", "SKU already exists")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	product := &models.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Cost:        req.Cost,
		Stock:       req.Stock,
		IsActive:    isActive,
	}

	if req.CategoryID != nil {
		catID, err := uuid.Parse(*req.CategoryID)
		if err == nil {
			product.CategoryID = &catID
		}
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to create product")
	}

	return product, nil
}

func (s *productService) GetByID(id uuid.UUID) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get product")
	}
	if product == nil {
		return nil, errors.ErrNotFound
	}
	return product, nil
}

func (s *productService) Update(id uuid.UUID, req *request.UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get product")
	}
	if product == nil {
		return nil, errors.ErrNotFound
	}

	if req.SKU != nil {
		existing, err := s.productRepo.GetBySKU(*req.SKU)
		if err != nil {
			return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to check SKU")
		}
		if existing != nil && existing.ID != id {
			return nil, errors.Wrap(errors.ErrDuplicateEntry, "SKU_EXISTS", "SKU already exists")
		}
		product.SKU = *req.SKU
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Cost != nil {
		product.Cost = *req.Cost
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.CategoryID != nil {
		if *req.CategoryID == "" {
			product.CategoryID = nil
		} else {
			catID, err := uuid.Parse(*req.CategoryID)
			if err == nil {
				product.CategoryID = &catID
			}
		}
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to update product")
	}

	return product, nil
}

func (s *productService) Delete(id uuid.UUID) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return errors.Wrap(err, "INTERNAL_ERROR", "Failed to get product")
	}
	if product == nil {
		return errors.ErrNotFound
	}

	product.IsActive = false
	if err := s.productRepo.Update(product); err != nil {
		return errors.Wrap(err, "INTERNAL_ERROR", "Failed to delete product")
	}

	return nil
}

func (s *productService) List(req *request.ListProductsRequest) ([]models.Product, int64, error) {
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

	var categoryID *uuid.UUID
	if req.CategoryID != "" {
		catID, err := uuid.Parse(req.CategoryID)
		if err == nil {
			categoryID = &catID
		}
	}

	products, total, err := s.productRepo.List(limit, offset, req.Search, categoryID, req.IsActive)
	if err != nil {
		return nil, 0, errors.Wrap(err, "INTERNAL_ERROR", "Failed to list products")
	}

	return products, total, nil
}

func (s *productService) GetLowStock(threshold int) ([]models.Product, error) {
	return s.productRepo.GetLowStock(threshold)
}
