package handler

import (
	"math"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"saleapp/internal/dto/request"
	"saleapp/internal/dto/response"
	"saleapp/internal/middleware"
	"saleapp/internal/models"
	"saleapp/internal/service"
	"saleapp/pkg/errors"
	pkgresponse "saleapp/pkg/response"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) List(c *gin.Context) {
	var req request.ListOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "query", Message: err.Error()},
		})
		return
	}

	orders, total, err := h.orderService.List(&req)
	if err != nil {
		pkgresponse.InternalError(c, "Failed to list orders")
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	orderResponses := make([]response.OrderResponse, len(orders))
	for i, o := range orders {
		orderResponses[i] = *response.NewOrderResponse(&o)
	}

	pkgresponse.SuccessWithMeta(c, gin.H{"orders": orderResponses}, &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "id", Message: "Invalid order ID"},
		})
		return
	}

	order, err := h.orderService.GetByID(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Order not found")
			return
		}
		pkgresponse.InternalError(c, "Failed to get order")
		return
	}

	pkgresponse.Success(c, response.NewOrderResponse(order))
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		pkgresponse.Unauthorized(c, "User not found in context")
		return
	}

	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	order, err := h.orderService.Create(userID, &req)
	if err != nil {
		if errors.Is(err, errors.ErrInvalidInput) {
			pkgresponse.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Resource not found")
			return
		}
		if errors.Is(err, errors.ErrInsufficientStock) {
			pkgresponse.Error(c, 422, "INSUFFICIENT_STOCK", err.Error())
			return
		}
		pkgresponse.InternalError(c, "Failed to create order")
		return
	}

	pkgresponse.Created(c, response.NewOrderResponse(order))
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "id", Message: "Invalid order ID"},
		})
		return
	}

	var req request.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	order, err := h.orderService.UpdateStatus(id, models.OrderStatus(req.Status))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Order not found")
			return
		}
		if errors.Is(err, errors.ErrInvalidStatus) {
			pkgresponse.Error(c, 422, "INVALID_STATUS", err.Error())
			return
		}
		pkgresponse.InternalError(c, "Failed to update order status")
		return
	}

	pkgresponse.Success(c, response.NewOrderResponse(order))
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "id", Message: "Invalid order ID"},
		})
		return
	}

	if err := h.orderService.Cancel(id); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Order not found")
			return
		}
		if errors.Is(err, errors.ErrInvalidStatus) {
			pkgresponse.Error(c, 422, "CANNOT_CANCEL", err.Error())
			return
		}
		pkgresponse.InternalError(c, "Failed to cancel order")
		return
	}

	pkgresponse.NoContent(c)
}
