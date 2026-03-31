package handler

import (
	"github.com/gin-gonic/gin"
	"saleapp/internal/dto/request"
	"saleapp/internal/middleware"
	"saleapp/internal/models"
	"saleapp/internal/service"
	"saleapp/pkg/errors"
	pkgresponse "saleapp/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	_, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errors.ErrUnauthorized) {
			pkgresponse.Unauthorized(c, "Invalid email or password")
			return
		}
		pkgresponse.InternalError(c, "Login failed")
		return
	}

	pkgresponse.Success(c, token)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgresponse.ValidationError(c, []pkgresponse.FieldError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	user, err := h.authService.Register(&models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.UserRole(req.Role),
	}, req.Password)
	if err != nil {
		if errors.Is(err, errors.ErrDuplicateEntry) {
			pkgresponse.Conflict(c, "Email already registered")
			return
		}
		pkgresponse.InternalError(c, "Registration failed")
		return
	}

	pkgresponse.Created(c, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		pkgresponse.Unauthorized(c, "User not found in context")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			pkgresponse.NotFound(c, "User not found")
			return
		}
		pkgresponse.InternalError(c, "Failed to get user")
		return
	}

	pkgresponse.Success(c, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	pkgresponse.Success(c, gin.H{"message": "Logged out successfully"})
}
