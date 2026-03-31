package handler

import (
	"math"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"saleapp/internal/dto/request"
	"saleapp/internal/dto/response"
	"saleapp/internal/service"
	"saleapp/pkg/errors"
	pkgresponse "saleapp/pkg/response"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

func (h *CustomerHandler) List(c *gin.Context) {
	var req request.ListCustomersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "query", Message: err.Error()},
		})
		return
	}

	customers, total, err := h.customerService.List(&req)
	if err != nil {
		pkgresponse.InternalError(c, "Failed to list customers")
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

	customerResponses := make([]response.CustomerResponse, len(customers))
	for i, cu := range customers {
		customerResponses[i] = *response.NewCustomerResponse(&cu)
	}

	pkgresponse.SuccessWithMeta(c, gin.H{"customers": customerResponses}, &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	})
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "id", Message: "Invalid customer ID"},
		})
		return
	}

	customer, err := h.customerService.GetByID(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Customer not found")
			return
		}
		pkgresponse.InternalError(c, "Failed to get customer")
		return
	}

	pkgresponse.Success(c, response.NewCustomerResponse(customer))
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req request.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	customer, err := h.customerService.Create(&req)
	if err != nil {
		if errors.Is(err, errors.ErrDuplicateEntry) {
			pkgresponse.Conflict(c, "Customer with this email already exists")
			return
		}
		pkgresponse.InternalError(c, "Failed to create customer")
		return
	}

	pkgresponse.Created(c, response.NewCustomerResponse(customer))
}

func (h *CustomerHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "id", Message: "Invalid customer ID"},
		})
		return
	}

	var req request.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	customer, err := h.customerService.Update(id, &req)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Customer not found")
			return
		}
		if errors.Is(err, errors.ErrDuplicateEntry) {
			pkgresponse.Conflict(c, "Customer with this email already exists")
			return
		}
		pkgresponse.InternalError(c, "Failed to update customer")
		return
	}

	pkgresponse.Success(c, response.NewCustomerResponse(customer))
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []response.FieldError{
			{Field: "id", Message: "Invalid customer ID"},
		})
		return
	}

	if err := h.customerService.Delete(id); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Customer not found")
			return
		}
		pkgresponse.InternalError(c, "Failed to delete customer")
		return
	}

	pkgresponse.NoContent(c)
}
