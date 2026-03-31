package service

import (
	"errors"
	"saleapp/internal/models"
	"saleapp/internal/repository"
	appErrors "saleapp/pkg/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (*models.User, string, error)
	Register(user *models.User, password string) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", appErrors.ErrUnauthorized
	}

	if !user.IsActive {
		return nil, "", appErrors.New("FORBIDDEN", "User account is inactive")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", appErrors.ErrUnauthorized
	}

	return user, "", nil
}

func (s *authService) Register(user *models.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return appErrors.Wrap(err, "INTERNAL_ERROR", "Failed to hash password")
	}

	user.PasswordHash = string(hashedPassword)
	user.Role = models.RoleCashier
	user.IsActive = true

	err = s.userRepo.Create(user)
	if err != nil {
		if appErrors.Is(err, appErrors.ErrDuplicateEntry) {
			return appErrors.New("CONFLICT", "Email already exists")
		}
		return err
	}

	return nil
}

func (s *authService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if appErrors.Is(err, appErrors.ErrNotFound) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *authService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return appErrors.ErrUnauthorized
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return appErrors.Wrap(err, "INTERNAL_ERROR", "Failed to hash password")
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(user)
}

// HashPassword is a helper function for external use (e.g., seeding)
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ValidatePassword checks if a password matches a hash
func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
