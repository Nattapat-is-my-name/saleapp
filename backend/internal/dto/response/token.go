package response

import (
	"saleapp/internal/models"
)

type TokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt string       `json:"expires_at"`
	User      *UserResponse `json:"user"`
}

func NewTokenResponse(token string, expiresAt string, user *models.User) *TokenResponse {
	return &TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
		},
	}
}
