package response

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"saleapp/internal/models"
)

type ProductResponse struct {
	ID          uuid.UUID       `json:"id"`
	SKU         string          `json:"sku"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Cost        decimal.Decimal `json:"cost"`
	Stock       int             `json:"stock"`
	CategoryID  *uuid.UUID      `json:"category_id,omitempty"`
	Category    *CategoryResponse `json:"category,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

type CategoryResponse struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

func NewProductResponse(p *models.Product) *ProductResponse {
	resp := &ProductResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Cost:        p.Cost,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if p.Category != nil {
		resp.Category = &CategoryResponse{
			ID:   p.Category.ID,
			Name: p.Category.Name,
			Slug: p.Category.Slug,
		}
	}
	return resp
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
}
