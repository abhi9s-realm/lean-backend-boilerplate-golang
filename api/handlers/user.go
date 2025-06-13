package handlers

import (
	"net/http"
	"strconv"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/services"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List all users with pagination
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, totalPages, totalItems, err := h.userService.ListUsers(c, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"current_page": page,
			"per_page":     limit,
			"total_items":  totalItems,
			"total_pages":  totalPages,
		},
	}
	utils.SuccessResponse(c, response, "Users fetched successfully")
}

// Get user by ID
func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	user, err := h.userService.GetUserByID(c, uint(id))
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
		}
		return
	}
	utils.SuccessResponse(c, user, "User fetched successfully")
}

// Create a new user
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest // Use DTO for request binding
	if err := c.ShouldBindJSON(&req); err != nil {
		// Use the custom ErrValidationFailed from the service layer for consistency, or define one in handlers
		utils.ErrorResponse(c, http.StatusBadRequest, services.ErrValidationFailed.Error()+": "+err.Error())
		return
	}

	// Map DTO to domain model
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	createdUser, err := h.userService.CreateUser(c, &user)
	if err != nil {
		if err == services.ErrUserEmailExists {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
		} else if err == services.ErrValidationFailed { // Should be caught by ShouldBindJSON generally
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		}
		return
	}
	utils.SuccessResponse(c, createdUser, "User created successfully")
}

// Update user
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var req UpdateUserRequest // Use DTO for request binding
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, services.ErrValidationFailed.Error()+": "+err.Error())
		return
	}

	// Map DTO to domain model for update.
	// Note: The service layer currently handles which fields are updatable (Name, Email).
	// This DTO helps ensure only these fields are considered from the request.
	userUpdate := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	updatedUser, err := h.userService.UpdateUser(c, uint(id), &userUpdate)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else if err == services.ErrEmailInUse {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
		} else if err == services.ErrValidationFailed { // Should be caught by ShouldBindJSON
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		}
		return
	}
	utils.SuccessResponse(c, updatedUser, "User updated successfully")
}

// Delete user
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	err = h.userService.DeleteUser(c, uint(id))
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		}
		return
	}
	utils.SuccessResponse(c, nil, "User deleted successfully")
}
