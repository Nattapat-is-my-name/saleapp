package response

import "github.com/shopspring/decimal"

type DashboardResponse struct {
	TodaySales   int64                 `json:"today_sales"`
	TodayOrders  int64                 `json:"today_orders"`
	TodayRevenue decimal.Decimal        `json:"today_revenue"`
	TopProducts  []TopProductResponse `json:"top_products"`
	LowStock     []ProductResponse    `json:"low_stock"`
}
