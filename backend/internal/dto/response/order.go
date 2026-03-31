package response

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"saleapp/internal/models"
)

type OrderItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	ProductID uuid.UUID       `json:"product_id"`
	Product   *ProductResponse `json:"product,omitempty"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
	Discount  decimal.Decimal `json:"discount"`
	Total     decimal.Decimal `json:"total"`
}

type CustomerResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Address   string    `json:"address"`
	Notes     string    `json:"notes"`
}

type UserResponse struct {
	ID        uuid.UUID       `json:"id"`
	Email     string          `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Role      models.UserRole `json:"role"`
}

type OrderResponse struct {
	ID            uuid.UUID           `json:"id"`
	OrderNumber   string              `json:"order_number"`
	Customer      *CustomerResponse   `json:"customer,omitempty"`
	User          UserResponse        `json:"user"`
	Status        models.OrderStatus  `json:"status"`
	Subtotal      decimal.Decimal     `json:"subtotal"`
	Tax           decimal.Decimal     `json:"tax"`
	Discount      decimal.Decimal     `json:"discount"`
	Total         decimal.Decimal     `json:"total"`
	PaymentMethod string              `json:"payment_method"`
	Notes         string              `json:"notes"`
	Items         []OrderItemResponse `json:"items"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
}

func NewOrderItemResponse(oi *models.OrderItem) *OrderItemResponse {
	resp := &OrderItemResponse{
		ID:        oi.ID,
		ProductID: oi.ProductID,
		Quantity:  oi.Quantity,
		UnitPrice: oi.UnitPrice,
		Discount:  oi.Discount,
		Total:     oi.Total,
	}
	if oi.Product.ID != uuid.Nil {
		resp.Product = NewProductResponse(&oi.Product)
	}
	return resp
}

func NewCustomerResponse(c *models.Customer) *CustomerResponse {
	if c == nil {
		return nil
	}
	return &CustomerResponse{
		ID:        c.ID,
		Email:     c.Email,
		Phone:     c.Phone,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Address:   c.Address,
		Notes:     c.Notes,
	}
}

func NewUserResponse(u *models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	}
}

func NewOrderResponse(o *models.Order) *OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = *NewOrderItemResponse(&item)
	}

	return &OrderResponse{
		ID:            o.ID,
		OrderNumber:   o.OrderNumber,
		Customer:      NewCustomerResponse(o.Customer),
		User:          NewUserResponse(&o.User),
		Status:        o.Status,
		Subtotal:      o.Subtotal,
		Tax:           o.Tax,
		Discount:      o.Discount,
		Total:         o.Total,
		PaymentMethod: o.PaymentMethod,
		Notes:         o.Notes,
		Items:         items,
		CreatedAt:     o.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int              `json:"total"`
}
