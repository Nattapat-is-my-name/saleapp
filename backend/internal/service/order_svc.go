package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"saleapp/internal/dto/request"
	"saleapp/internal/dto/response"
	"saleapp/internal/models"
	"saleapp/internal/repository"
	"saleapp/pkg/errors"
)

type OrderService interface {
	Create(userID uuid.UUID, req *request.CreateOrderRequest) (*models.Order, error)
	GetByID(id uuid.UUID) (*models.Order, error)
	UpdateStatus(id uuid.UUID, status models.OrderStatus) (*models.Order, error)
	Cancel(id uuid.UUID) error
	List(req *request.ListOrdersRequest) ([]models.Order, int64, error)
	GetSalesSummary(startDate, endDate time.Time) (*response.SalesSummaryResponse, error)
	GetTopSellingProducts(startDate, endDate time.Time, limit int) ([]response.TopProductResponse, error)
	GetLowStockProducts(threshold int) ([]models.Product, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	customerRepo repository.CustomerRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, customerRepo repository.CustomerRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		customerRepo: customerRepo,
	}
}

func (s *orderService) Create(userID uuid.UUID, req *request.CreateOrderRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.Wrap(errors.ErrInvalidInput, "NO_ITEMS", "Order must have at least one item")
	}

	orderNumber := fmt.Sprintf("ORD-%d-%04d", time.Now().Unix(), time.Now().Nanosecond()%10000)

	var customerID *uuid.UUID
	if req.CustomerID != nil && *req.CustomerID != "" {
		cID, err := uuid.Parse(*req.CustomerID)
		if err != nil {
			return nil, errors.Wrap(errors.ErrInvalidInput, "INVALID_CUSTOMER", "Invalid customer ID")
		}
		customer, err := s.customerRepo.GetByID(cID)
		if err != nil {
			return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get customer")
		}
		if customer == nil {
			return nil, errors.Wrap(errors.ErrNotFound, "CUSTOMER_NOT_FOUND", "Customer not found")
		}
		customerID = &cID
	}

	order := &models.Order{
		OrderNumber:   orderNumber,
		CustomerID:    customerID,
		UserID:        userID,
		Status:        models.StatusPending,
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
		Items:         make([]models.OrderItem, 0, len(req.Items)),
	}

	var subtotal decimal.Decimal
	taxRate := decimal.NewFromFloat(0.08) // 8% tax

	for _, itemReq := range req.Items {
		productID, err := uuid.Parse(itemReq.ProductID)
		if err != nil {
			return nil, errors.Wrap(errors.ErrInvalidInput, "INVALID_PRODUCT", "Invalid product ID")
		}

		product, err := s.productRepo.GetByID(productID)
		if err != nil {
			return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get product")
		}
		if product == nil {
			return nil, errors.Wrap(errors.ErrNotFound, "PRODUCT_NOT_FOUND", "Product not found")
		}

		if !product.IsInStock() {
			return nil, errors.Wrapf(errors.ErrInsufficientStock, "INSUFFICIENT_STOCK", "Product %s is out of stock", product.Name)
		}

		if product.Stock < itemReq.Quantity {
			return nil, errors.Wrapf(errors.ErrInsufficientStock, "INSUFFICIENT_STOCK", 
				"Product %s has insufficient stock (available: %d, requested: %d)", product.Name, product.Stock, itemReq.Quantity)
		}

		itemTotal := product.Price.Mul(decimal.NewFromInt(int64(itemReq.Quantity))).Sub(itemReq.Discount)

		orderItem := models.OrderItem{
			ProductID: productID,
			Quantity:  itemReq.Quantity,
			UnitPrice: product.Price,
			Discount:  itemReq.Discount,
			Total:     itemTotal,
		}

		order.Items = append(order.Items, orderItem)
		subtotal = subtotal.Add(itemTotal)

		if err := s.productRepo.UpdateStock(productID, -itemReq.Quantity); err != nil {
			return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to update stock")
		}
	}

	order.Subtotal = subtotal
	order.Tax = subtotal.Mul(taxRate).Round(2)
	order.Discount = decimal.Zero
	order.Total = order.Subtotal.Add(order.Tax).Sub(order.Discount)

	if err := s.orderRepo.Create(order); err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to create order")
	}

	order.Status = models.StatusCompleted
	s.orderRepo.Update(order)

	createdOrder, err := s.orderRepo.GetByID(order.ID)
	if err != nil {
		return order, nil
	}
	return createdOrder, nil
}

func (s *orderService) GetByID(id uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get order")
	}
	if order == nil {
		return nil, errors.ErrNotFound
	}
	return order, nil
}

func (s *orderService) UpdateStatus(id uuid.UUID, status models.OrderStatus) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get order")
	}
	if order == nil {
		return nil, errors.ErrNotFound
	}

	validStatuses := map[models.OrderStatus][]models.OrderStatus{
		models.StatusPending:   {models.StatusCompleted, models.StatusCancelled},
		models.StatusCompleted: {models.StatusRefunded},
	}

	allowed, exists := validStatuses[order.Status]
	if !exists {
		return nil, errors.Wrap(errors.ErrInvalidStatus, "INVALID_TRANSITION", 
			fmt.Sprintf("Cannot change status from %s", order.Status))
	}

	valid := false
	for _, s := range allowed {
		if s == status {
			valid = true
			break
		}
	}

	if !valid {
		return nil, errors.Wrap(errors.ErrInvalidStatus, "INVALID_TRANSITION",
			fmt.Sprintf("Cannot transition from %s to %s", order.Status, status))
	}

	if status == models.StatusCancelled || status == models.StatusRefunded {
		for _, item := range order.Items {
			s.productRepo.UpdateStock(item.ProductID, item.Quantity)
		}
	}

	order.Status = status
	if err := s.orderRepo.Update(order); err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to update order")
	}

	return order, nil
}

func (s *orderService) Cancel(id uuid.UUID) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.Wrap(err, "INTERNAL_ERROR", "Failed to get order")
	}
	if order == nil {
		return errors.ErrNotFound
	}

	if order.Status == models.StatusCancelled {
		return errors.ErrOrderCancelled
	}

	if order.Status != models.StatusPending && order.Status != models.StatusCompleted {
		return errors.Wrap(errors.ErrInvalidStatus, "CANNOT_CANCEL", "Order cannot be cancelled in current status")
	}

	for _, item := range order.Items {
		if err := s.productRepo.UpdateStock(item.ProductID, item.Quantity); err != nil {
			return errors.Wrap(err, "INTERNAL_ERROR", "Failed to restore stock")
		}
	}

	order.Status = models.StatusCancelled
	if err := s.orderRepo.Update(order); err != nil {
		return errors.Wrap(err, "INTERNAL_ERROR", "Failed to cancel order")
	}

	return nil
}

func (s *orderService) List(req *request.ListOrdersRequest) ([]models.Order, int64, error) {
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

	var customerID *uuid.UUID
	if req.CustomerID != "" {
		cID, err := uuid.Parse(req.CustomerID)
		if err == nil {
			customerID = &cID
		}
	}

	var startDate, endDate *time.Time
	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			startDate = &t
		}
	}
	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			end := t.Add(24*time.Hour - time.Second)
			endDate = &end
		}
	}

	orders, total, err := s.orderRepo.List(limit, offset, customerID, req.Status, startDate, endDate)
	if err != nil {
		return nil, 0, errors.Wrap(err, "INTERNAL_ERROR", "Failed to list orders")
	}

	return orders, total, nil
}

func (s *orderService) GetSalesSummary(startDate, endDate time.Time) (*response.SalesSummaryResponse, error) {
	orders, err := s.orderRepo.GetSalesByDate(startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get sales")
	}

	var totalRevenue decimal.Decimal
	var totalOrders int64
	orderCounts := make(map[string]int)

	for _, order := range orders {
		if order.Status == models.StatusCompleted {
			totalRevenue = totalRevenue.Add(order.Total)
			totalOrders++
		}
		orderCounts[order.OrderNumber]++
	}

	return &response.SalesSummaryResponse{
		TotalRevenue:  totalRevenue,
		TotalOrders:   totalOrders,
		AverageOrder:  totalRevenue.Div(decimal.NewFromInt(totalOrders)).Round(2),
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
	}, nil
}

func (s *orderService) GetTopSellingProducts(startDate, endDate time.Time, limit int) ([]response.TopProductResponse, error) {
	results, err := s.orderRepo.GetTopSellingProducts(startDate, endDate, limit)
	if err != nil {
		return nil, errors.Wrap(err, "INTERNAL_ERROR", "Failed to get top selling products")
	}

	products := make([]response.TopProductResponse, len(results))
	for i, r := range results {
		product, _ := s.productRepo.GetByID(r.ProductID)
		if product != nil {
			products[i] = response.TopProductResponse{
				Product:  *response.NewProductResponse(product),
				Quantity: r.TotalSold,
			}
		}
	}

	return products, nil
}

func (s *orderService) GetLowStockProducts(threshold int) ([]models.Product, error) {
	return s.productRepo.GetLowStock(threshold)
}
