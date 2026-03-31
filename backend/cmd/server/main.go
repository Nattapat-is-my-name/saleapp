package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"saleapp/internal/config"
	"saleapp/internal/handler"
	"saleapp/internal/middleware"
	"saleapp/internal/models"
	"saleapp/internal/repository"
	"saleapp/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.Server.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Connect to database
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Auto migrate
	err = db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Product{},
		&models.Category{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// Seed default admin user if not exists
	seedAdminUser(db)

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	authService := service.NewAuthService(userRepo)
	productService := service.NewProductService(productRepo)
	customerService := service.NewCustomerService(customerRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, userRepo, customerRepo)

	// JWT Middleware
	jwtMw := middleware.NewJWTMiddleware(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, jwtMw)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customerService)
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup Gin
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", jwtMw.AuthRequired(), authHandler.Logout)
			auth.GET("/me", jwtMw.AuthRequired(), authHandler.Me)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(jwtMw.AuthRequired())
		{
			// Products
			products := protected.Group("/products")
			{
				products.GET("", productHandler.List)
				products.GET("/:id", productHandler.Get)
				products.POST("", productHandler.Create)
				products.PUT("/:id", productHandler.Update)
				products.DELETE("/:id", productHandler.Delete)
				products.GET("/low-stock", productHandler.GetLowStock)
			}

			// Customers
			customers := protected.Group("/customers")
			{
				customers.GET("", customerHandler.List)
				customers.GET("/:id", customerHandler.Get)
				customers.POST("", customerHandler.Create)
				customers.PUT("/:id", customerHandler.Update)
				customers.DELETE("/:id", customerHandler.Delete)
			}

			// Orders
			orders := protected.Group("/orders")
			{
				orders.GET("", orderHandler.List)
				orders.GET("/:id", orderHandler.Get)
				orders.POST("", orderHandler.Create)
				orders.PUT("/:id/status", orderHandler.UpdateStatus)
				orders.DELETE("/:id", orderHandler.Cancel)
			}

			// Reports (admin/manager only)
			reports := protected.Group("/reports")
			reports.Use(jwtMw.RoleRequired(models.RoleAdmin, models.RoleManager))
			{
				reports.GET("/sales", orderHandler.GetSalesReport)
				reports.GET("/products/top", orderHandler.GetTopProducts)
			}
		}
	}

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Info().Msgf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.Server.Env == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Warn)
	}

	dsn := cfg.Database.DSN()
	return gorm.Open(postgres.Open(dsn), gormConfig)
}

func seedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@saleapp.local").Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, err := service.HashPassword("admin123")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to hash admin password")
		return
	}

	admin := &models.User{
		Email:        "admin@saleapp.local",
		PasswordHash: hashedPassword,
		FirstName:    "Admin",
		LastName:     "User",
		Role:         models.RoleAdmin,
		IsActive:     true,
	}

	if err := db.Create(admin).Error; err != nil {
		log.Warn().Err(err).Msg("Failed to create admin user")
		return
	}

	log.Info().Msg("Default admin user created: admin@saleapp.local / admin123")
}
