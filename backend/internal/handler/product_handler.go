package handler

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"saleapp/internal/dto/request"
	"saleapp/internal/dto/response"
	"saleapp/internal/service"
	"saleapp/pkg/errors"
	pkgresponse "saleapp/pkg/response"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) List(c *gin.Context) {
	var req request.ListProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "query", Message: err.Error()},
		})
		return
	}

	products, total, err := h.productService.List(&req)
	if err != nil {
		pkgresponse.InternalError(c, "Failed to list products")
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

	productResponses := make([]response.ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = *response.NewProductResponse(&p)
	}

	pkgresponse.SuccessWithMeta(c, gin.H{"products": productResponses}, &pkgresponse.Meta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "id", Message: "Invalid product ID"},
		})
		return
	}

	product, err := h.productService.GetByID(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Product not found")
			return
		}
		pkgresponse.InternalError(c, "Failed to get product")
		return
	}

	pkgresponse.Success(c, response.NewProductResponse(product))
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	product, err := h.productService.Create(&req)
	if err != nil {
		if errors.Is(err, errors.ErrDuplicateEntry) {
			pkgresponse.Conflict(c, "SKU already exists")
			return
		}
		pkgresponse.Error(c, http.StatusUnprocessableEntity, "CREATE_FAILED", err.Error())
		return
	}

	pkgresponse.Created(c, response.NewProductResponse(product))
}

func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "id", Message: "Invalid product ID"},
		})
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	product, err := h.productService.Update(id, &req)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Product not found")
			return
		}
		if errors.Is(err, errors.ErrDuplicateEntry) {
			pkgresponse.Conflict(c, "SKU already exists")
			return
		}
		pkgresponse.InternalError(c, "Failed to update product")
		return
	}

	pkgresponse.Success(c, response.NewProductResponse(product))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "id", Message: "Invalid product ID"},
		})
		return
	}

	if err := h.productService.Delete(id); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "Product not found")
			return
		}
		pkgresponse.InternalError(c, "Failed to delete product")
		return
	}

	pkgresponse.NoContent(c)
}
