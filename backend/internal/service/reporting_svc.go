package service

import (
	"time"

	"github.com/shopspring/decimal"
	"saleapp/internal/dto/response"
	"saleapp/internal/repository"
)

type ReportingService interface {
	GetSalesSummary(startDate, endDate time.Time) (*response.SalesSummaryResponse, error)
	GetTopSellingProducts(startDate, endDate time.Time, limit int) ([]response.TopProductResponse, error)
	GetLowStockProducts(threshold int) (*response.LowStockResponse, error)
	GetDashboard() (*response.DashboardResponse, error)
}

type reportingService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewReportingService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) ReportingService {
	return &reportingService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *reportingService) GetSalesSummary(startDate, endDate time.Time) (*response.SalesSummaryResponse, error) {
	orders, err := s.orderRepo.GetSalesByDate(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var totalRevenue decimal.Decimal
	var totalOrders int64

	for _, order := range orders {
		if order.Status == "completed" {
			totalRevenue = totalRevenue.Add(order.Total)
			totalOrders++
		}
	}

	var avgOrder decimal.Decimal
	if totalOrders > 0 {
		avgOrder = totalRevenue.Div(decimal.NewFromInt(totalOrders))
	}

	return &response.SalesSummaryResponse{
		TotalRevenue:  totalRevenue,
		TotalOrders:   totalOrders,
		AverageOrder:  avgOrder,
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
	}, nil
}

func (s *reportingService) GetTopSellingProducts(startDate, endDate time.Time, limit int) ([]response.TopProductResponse, error) {
	results, err := s.orderRepo.GetTopSellingProducts(startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	products := make([]response.TopProductResponse, 0, len(results))
	for _, r := range results {
		product, err := s.productRepo.GetByID(r.ProductID)
		if err != nil {
			continue
		}
		if product != nil {
			products = append(products, response.TopProductResponse{
				Product:  *response.NewProductResponse(product),
				Quantity: r.TotalSold,
			})
		}
	}

	return products, nil
}

func (s *reportingService) GetLowStockProducts(threshold int) (*response.LowStockResponse, error) {
	products, err := s.productRepo.GetLowStock(threshold)
	if err != nil {
		return nil, err
	}

	productResponses := make([]response.ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = *response.NewProductResponse(&p)
	}

	return &response.LowStockResponse{
		Products:  productResponses,
		Threshold: threshold,
	}, nil
}

func (s *reportingService) GetDashboard() (*response.DashboardResponse, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Second)

	// Get today's completed orders
	orders, err := s.orderRepo.GetSalesByDate(startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	var todayRevenue decimal.Decimal
	var todayOrders int64
	for _, order := range orders {
		if order.Status == "completed" {
			todayRevenue = todayRevenue.Add(order.Total)
			todayOrders++
		}
	}

	// Get top 5 selling products for today
	topResults, err := s.orderRepo.GetTopSellingProducts(startOfDay, endOfDay, 5)
	if err != nil {
		return nil, err
	}

	topProducts := make([]response.TopProductResponse, 0, len(topResults))
	for _, r := range topResults {
		product, err := s.productRepo.GetByID(r.ProductID)
		if err != nil || product == nil {
			continue
		}
		topProducts = append(topProducts, response.TopProductResponse{
			Product:      *response.NewProductResponse(product),
			Quantity: r.TotalSold,
		})
	}

	// Get low stock products (threshold 10)
	lowStockProducts, err := s.productRepo.GetLowStock(10)
	if err != nil {
		return nil, err
	}

	lowStock := make([]response.ProductResponse, len(lowStockProducts))
	for i, p := range lowStockProducts {
		lowStock[i] = *response.NewProductResponse(&p)
	}

	return &response.DashboardResponse{
		TodaySales:   todayOrders,
		TodayOrders:  todayOrders,
		TodayRevenue: todayRevenue,
		TopProducts:  topProducts,
		LowStock:     lowStock,
	}, nil
}
