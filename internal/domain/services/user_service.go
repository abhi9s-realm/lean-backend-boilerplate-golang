package services

import (
	"errors"
	"math"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Required for gorm.ErrRecordNotFound, potentially define custom errors later
)

// Define custom errors (good practice, can be moved to a dedicated errors package)
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserEmailExists  = errors.New("user with this email already exists")
	ErrEmailInUse       = errors.New("email already in use by another user")
	ErrValidationFailed = errors.New("validation failed")
)

type userServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) ListUsers(c *gin.Context, page, limit int) ([]models.User, int, int64, error) {
	users, total, err := s.userRepo.List(c, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}
	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}
	return users, totalPages, total, nil
}

func (s *userServiceImpl) GetUserByID(c *gin.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(c, id)
	if err != nil {
		return nil, err // Pass through repository error
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userServiceImpl) CreateUser(c *gin.Context, user *models.User) (*models.User, error) {
	// Business logic: Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(c, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && existingUser != nil { // Check if it is not simply a "not found" error
		return nil, err // Database error
	}
	if existingUser != nil {
		return nil, ErrUserEmailExists
	}

	if err := s.userRepo.Create(c, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) UpdateUser(c *gin.Context, id uint, userUpdate *models.User) (*models.User, error) {
	existingUser, err := s.userRepo.GetByID(c, id)
	if err != nil {
		return nil, err // Database error
	}
	if existingUser == nil {
		return nil, ErrUserNotFound
	}

	// Business logic: Check if email is being changed and if it already exists for another user
	if userUpdate.Email != "" && userUpdate.Email != existingUser.Email {
		collidingUser, err := s.userRepo.GetByEmail(c, userUpdate.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && collidingUser != nil {
			return nil, err // Database error
		}
		if collidingUser != nil && collidingUser.ID != existingUser.ID {
			return nil, ErrEmailInUse
		}
		existingUser.Email = userUpdate.Email
	}

	if userUpdate.Name != "" {
		existingUser.Name = userUpdate.Name
	}
	// Potentially update other fields as needed

	if err := s.userRepo.Update(c, existingUser); err != nil {
		return nil, err
	}
	return existingUser, nil
}

func (s *userServiceImpl) DeleteUser(c *gin.Context, id uint) error {
	user, err := s.userRepo.GetByID(c, id)
	if err != nil {
		return err // Database error
	}
	if user == nil {
		return ErrUserNotFound
	}
	return s.userRepo.Delete(c, id)
}
