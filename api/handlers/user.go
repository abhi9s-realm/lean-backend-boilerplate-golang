package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// List all users with pagination
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	if err := h.db.Model(&models.User{}).Count(&total).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to count users")
		return
	}

	if err := h.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	response := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	}

	utils.SuccessResponse(c, response, "Users fetched successfully")
}

// Get user by ID
func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, user, "User fetched successfully")
}

// Create a new user
func (h *UserHandler) Create(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// Check if user with email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists")
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	utils.SuccessResponse(c, user, "User created successfully")
}

// Update user
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var existingUser models.User
	if err := h.db.First(&existingUser, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// Check if email is being changed and if it already exists
	if updateData.Email != existingUser.Email {
		var emailCheck models.User
		if err := h.db.Where("email = ?", updateData.Email).First(&emailCheck).Error; err == nil {
			utils.ErrorResponse(c, http.StatusConflict, "Email already in use")
			return
		}
	}

	// Update allowed fields
	existingUser.Name = updateData.Name
	existingUser.Email = updateData.Email

	if err := h.db.Save(&existingUser).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}

	utils.SuccessResponse(c, existingUser, "User updated successfully")
}

// Delete user
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.SuccessResponse(c, nil, "User deleted successfully")
}
