package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"saleapp/internal/config"
	"saleapp/internal/handler"
	"saleapp/internal/middleware"
	"saleapp/internal/models"
	"saleapp/internal/repository"
	"saleapp/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Initialize zerolog
	zlog := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		// Try with default config for development
		cfg = &config.Config{
			Server: config.ServerConfig{
				Host: "0.0.0.0",
				Port: "8080",
				Env:  "development",
			},
			Database: config.DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				Name:     "saleapp",
				User:     "postgres",
				Password: "postgres",
				SSLMode:  "disable",
			},
			JWT: config.JWTConfig{
				Secret:      "your-secret-key-change-in-production-min-32-chars",
				ExpiryHours: 24,
			},
			Log: config.LogConfig{
				Level: "debug",
			},
		}
		zlog.Warn().Err(err).Msg("Failed to load config, using defaults")
	}

	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Override database config from env vars (for Docker)
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		cfg.Database.Port = port
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if sslmode := os.Getenv("DB_SSL_MODE"); sslmode != "" {
		cfg.Database.SSLMode = sslmode
	}

	// Initialize database
	db, err := initDatabase(cfg, &zlog)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Failed to initialize database")
	}

	// JWT Middleware
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	productService := service.NewProductService(productRepo)
	customerService := service.NewCustomerService(customerRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, customerRepo)
	reportingService := service.NewReportingService(orderRepo, productRepo)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo)
	paymentService.InitStripe()

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, jwtMiddleware)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customerService)
	orderHandler := handler.NewOrderHandler(orderService)
	reportingHandler := handler.NewReportingHandler(reportingService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Rate limiter
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Setup router
	router := gin.New()

	// Middleware
	router.Use(middleware.Recovery(zlog))
	router.Use(middleware.Logger(zlog))
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	router.HEAD("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Stripe webhook (must be before auth - Stripe calls this directly)
		v1.POST("/payments/webhook", paymentHandler.HandleWebhook)

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", jwtMiddleware.AuthRequired(), authHandler.Me)
		}

		// Payment routes (protected)
		payments := v1.Group("/payments")
		payments.Use(jwtMiddleware.AuthRequired())
		{
			payments.POST("/intent", paymentHandler.CreatePaymentIntent)
			payments.GET("/:orderId", paymentHandler.GetPaymentStatus)
		}

		// Products routes (protected)
		products := v1.Group("/products")
		products.Use(jwtMiddleware.AuthRequired())
		products.Use(middleware.RateLimit(rateLimiter))
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.GetByID)
			products.POST("", jwtMiddleware.RoleRequired(models.RoleAdmin, models.RoleManager), productHandler.Create)
			products.PUT("/:id", jwtMiddleware.RoleRequired(models.RoleAdmin, models.RoleManager), productHandler.Update)
			products.DELETE("/:id", jwtMiddleware.RoleRequired(models.RoleAdmin), productHandler.Delete)
		}

		// Customers routes (protected)
		customers := v1.Group("/customers")
		customers.Use(jwtMiddleware.AuthRequired())
		customers.Use(middleware.RateLimit(rateLimiter))
		{
			customers.GET("", customerHandler.List)
			customers.GET("/:id", customerHandler.GetByID)
			customers.POST("", jwtMiddleware.RoleRequired(models.RoleAdmin, models.RoleManager), customerHandler.Create)
			customers.PUT("/:id", jwtMiddleware.RoleRequired(models.RoleAdmin, models.RoleManager), customerHandler.Update)
			customers.DELETE("/:id", jwtMiddleware.RoleRequired(models.RoleAdmin), customerHandler.Delete)
		}

		// Orders routes (protected)
		orders := v1.Group("/orders")
		orders.Use(jwtMiddleware.AuthRequired())
		orders.Use(middleware.RateLimit(rateLimiter))
		{
			orders.GET("", orderHandler.List)
			orders.GET("/:id", orderHandler.GetByID)
			orders.POST("", orderHandler.Create)
			orders.PUT("/:id/status", jwtMiddleware.RoleRequired(models.RoleAdmin, models.RoleManager), orderHandler.UpdateStatus)
			orders.DELETE("/:id", jwtMiddleware.RoleRequired(models.RoleAdmin, models.RoleManager), orderHandler.Cancel)
		}

		// Reports routes (protected)
		reports := v1.Group("/reports")
		reports.Use(jwtMiddleware.AuthRequired())
		reports.Use(middleware.RateLimit(rateLimiter))
		{
			reports.GET("/sales", reportingHandler.GetSalesSummary)
			reports.GET("/products/top", reportingHandler.GetTopSellingProducts)
			reports.GET("/inventory/low", reportingHandler.GetLowStockProducts)
			reports.GET("/dashboard", reportingHandler.GetDashboard)
		}
	}

	// Create server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		zlog.Info().Str("address", addr).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zlog.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zlog.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	zlog.Info().Msg("Server exited")
}

func initDatabase(cfg *config.Config, zlog *zerolog.Logger) (*gorm.DB, error) {
	// Check if we should use SQLite for development
	dsn := cfg.Database.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		// Fallback to SQLite for local development
		zlog.Warn().Msg("PostgreSQL connection failed, using SQLite for development")
		dsn = "file::memory:?cache=shared"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Customer{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.ProcessedEvent{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	zlog.Info().Msg("Database connected and migrated")
	return db, nil
}
