package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"saleapp/internal/service"
	pkgresponse "saleapp/pkg/response"
)

type ReportingHandler struct {
	reportingService service.ReportingService
}

func NewReportingHandler(reportingService service.ReportingService) *ReportingHandler {
	return &ReportingHandler{
		reportingService: reportingService,
	}
}

func (h *ReportingHandler) GetSalesSummary(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" {
		startDateStr = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "start_date", Message: "Invalid date format (use YYYY-MM-DD)"},
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "end_date", Message: "Invalid date format (use YYYY-MM-DD)"},
		})
		return
	}
	endDate = endDate.Add(24*time.Hour - time.Second)

	summary, err := h.reportingService.GetSalesSummary(startDate, endDate)
	if err != nil {
		pkgresponse.InternalError(c, "Failed to get sales summary")
		return
	}

	pkgresponse.Success(c, summary)
}

func (h *ReportingHandler) GetTopSellingProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" {
		startDateStr = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "start_date", Message: "Invalid date format (use YYYY-MM-DD)"},
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "end_date", Message: "Invalid date format (use YYYY-MM-DD)"},
		})
		return
	}
	endDate = endDate.Add(24*time.Hour - time.Second)

	products, err := h.reportingService.GetTopSellingProducts(startDate, endDate, limit)
	if err != nil {
		pkgresponse.InternalError(c, "Failed to get top selling products")
		return
	}

	pkgresponse.Success(c, gin.H{"products": products})
}

func (h *ReportingHandler) GetLowStockProducts(c *gin.Context) {
	thresholdStr := c.DefaultQuery("threshold", "10")
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil || threshold < 1 {
		threshold = 10
	}

	lowStock, err := h.reportingService.GetLowStockProducts(threshold)
	if err != nil {
		pkgresponse.InternalError(c, "Failed to get low stock products")
		return
	}

	pkgresponse.Success(c, lowStock)
}
