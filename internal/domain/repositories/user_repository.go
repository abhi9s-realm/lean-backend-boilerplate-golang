package repositories

import (
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	List(c *gin.Context, page, limit int) ([]models.User, int64, error)
	GetByID(c *gin.Context, id uint) (*models.User, error)
	GetByEmail(c *gin.Context, email string) (*models.User, error)
	Create(c *gin.Context, user *models.User) error
	Update(c *gin.Context, user *models.User) error
	Delete(c *gin.Context, id uint) error
}
