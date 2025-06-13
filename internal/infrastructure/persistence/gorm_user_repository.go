package persistence

import (
	"errors"
	"math"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repositories.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) List(c *gin.Context, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	offset := (page - 1) * limit

	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *GormUserRepository) GetByID(c *gin.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom domain error e.g. ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByEmail(c *gin.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom domain error e.g. ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Create(c *gin.Context, user *models.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) Update(c *gin.Context, user *models.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) Delete(c *gin.Context, id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
