package response

import (
	"github.com/shopspring/decimal"
)

type SalesSummaryResponse struct {
	TotalRevenue decimal.Decimal `json:"total_revenue"`
	TotalOrders  int64           `json:"total_orders"`
	AverageOrder decimal.Decimal `json:"average_order"`
	StartDate    string          `json:"start_date"`
	EndDate      string          `json:"end_date"`
}

type TopProductResponse struct {
	Product  ProductResponse `json:"product"`
	Quantity int64           `json:"quantity_sold"`
}

type LowStockResponse struct {
	Products []ProductResponse `json:"products"`
	Threshold int              `json:"threshold"`
}
