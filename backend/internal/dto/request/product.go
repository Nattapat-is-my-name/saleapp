package request

import "github.com/shopspring/decimal"

type CreateProductRequest struct {
	SKU         string           `json:"sku" binding:"required,max=50"`
	Name        string           `json:"name" binding:"required,max=255"`
	Description string           `json:"description"`
	Price       decimal.Decimal `json:"price" binding:"required"`
	Cost        decimal.Decimal `json:"cost"`
	Stock       int              `json:"stock" binding:"min=0"`
	CategoryID  *string          `json:"category_id"`
	IsActive    *bool            `json:"is_active"`
}

type UpdateProductRequest struct {
	SKU         *string           `json:"sku" binding:"omitempty,max=50"`
	Name        *string           `json:"name" binding:"omitempty,max=255"`
	Description *string           `json:"description"`
	Price       *decimal.Decimal `json:"price"`
	Cost        *decimal.Decimal `json:"cost"`
	Stock       *int              `json:"stock" binding:"omitempty,min=0"`
	CategoryID  *string           `json:"category_id"`
	IsActive    *bool             `json:"is_active"`
}

type ListProductsRequest struct {
	Page       int    `form:"page" binding:"min=0"`
	Limit      int    `form:"limit" binding:"min=0,max=100"`
	Search     string `form:"search"`
	CategoryID string `form:"category_id"`
	IsActive   *bool  `form:"is_active"`
}
