package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/api/handlers"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/models"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/pkg/utils"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUserCRUD(t *testing.T) {
	// Setup with database since user tests need DB access
	r, db := testutils.SetupTestRouter(true)
	if db == nil {
		t.Fatal("Database connection required for user tests")
	}
	
	// Clean up database before starting tests
	testutils.CleanupDatabase(db)

	userHandler := handlers.NewUserHandler(db)
	r.POST("/api/users", userHandler.Create)
	r.GET("/api/users", userHandler.List)
	r.GET("/api/users/:id", userHandler.Get)
	r.PUT("/api/users/:id", userHandler.Update)
	r.DELETE("/api/users/:id", userHandler.Delete)

	t.Run("Create User", func(t *testing.T) {
		userData := map[string]interface{}{
			"name":  "Test User",
			"email": "test@example.com",
		}
		jsonData, _ := json.Marshal(userData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)

		// Verify user exists in database
		var user models.User
		err = db.Where("email = ?", userData["email"]).First(&user).Error
		assert.NoError(t, err, "User should exist in database after creation")
	})

	t.Run("List Users After Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		// Check the response structure
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")

		users, ok := data["users"].([]interface{})
		assert.True(t, ok, "Response should contain users array")
		assert.NotEmpty(t, users, "Users list should not be empty after creating a user")
	})

	t.Run("Get User", func(t *testing.T) {
		// First create a user
		user := models.User{Name: "Get Test User", Email: "get@example.com"}
		db.Create(&user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/users/%d", user.ID), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Update User", func(t *testing.T) {
		// First create a user
		user := models.User{Name: "Update Test User", Email: "update@example.com"}
		db.Create(&user)

		updateData := map[string]interface{}{
			"name":  "Updated User",
			"email": "updated@example.com",
		}
		jsonData, _ := json.Marshal(updateData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/users/%d", user.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Delete User", func(t *testing.T) {
		// First create a user
		user := models.User{Name: "Delete Test User", Email: "delete@example.com"}
		db.Create(&user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/users/%d", user.ID), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		// Verify user is deleted
		var deletedUser models.User
		err = db.First(&deletedUser, user.ID).Error
		assert.Error(t, err) // Should not find the user
	})
}
