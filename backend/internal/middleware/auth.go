package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"saleapp/internal/models"
	"saleapp/pkg/errors"
	"saleapp/pkg/response"
)

type Claims struct {
	UserID uuid.UUID   `json:"user_id"`
	Email string      `json:"email"`
	Role  models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type JWTMiddleware struct {
	secret      []byte
	expiryHours int
}

func NewJWTMiddleware(secret string, expiryHours int) *JWTMiddleware {
	return &JWTMiddleware{
		secret:      []byte(secret),
		expiryHours: expiryHours,
	}
}

func (m *JWTMiddleware) GenerateToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(m.expiryHours) * time.Hour)
	
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "saleapp",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (m *JWTMiddleware) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrUnauthorized
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrUnauthorized
}

func (m *JWTMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := m.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func (m *JWTMiddleware) RoleRequired(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "User role not found in context")
			c.Abort()
			return
		}

		userRole := roleVal.(models.UserRole)
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	return userID.(uuid.UUID), true
}

func GetUserRole(c *gin.Context) (models.UserRole, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return role.(models.UserRole), true
}
