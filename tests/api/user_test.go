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
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/internal/domain/services" // Needed for error comparison
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/pkg/utils"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/tests/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserCRUD(t *testing.T) {
	r, userService, db := testutils.SetupTestRouter(true)
	if db == nil {
		t.Fatal("Database connection required for user tests")
	}
	if userService == nil {
		t.Fatal("UserService required for user tests")
	}
	
	testutils.CleanupDatabase(db)

	userHandler := handlers.NewUserHandler(userService)

	userRoutes := r.Group("/api/users")
	{
		userRoutes.POST("", userHandler.Create)
		userRoutes.GET("", userHandler.List)
		userRoutes.GET("/:id", userHandler.Get)
		userRoutes.PUT("/:id", userHandler.Update)
		userRoutes.DELETE("/:id", userHandler.Delete)
	}


	t.Run("Create User", func(t *testing.T) {
		defer testutils.CleanupDatabase(db)
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

        responseData, ok := response.Data.(map[string]interface{})
        assert.True(t, ok)
        assert.Equal(t, userData["name"], responseData["name"])
        assert.Equal(t, userData["email"], responseData["email"])
        assert.NotNil(t, responseData["id"])

		var user models.User
		err = db.Where("email = ?", userData["email"]).First(&user).Error
		assert.NoError(t, err, "User should exist in database after creation")
		assert.Equal(t, userData["name"], user.Name)
	})

	t.Run("List Users After Create", func(t *testing.T) {
		defer testutils.CleanupDatabase(db)
		initialUser := models.User{Name: "List Test User", Email: "list@example.com"}
		db.Create(&initialUser)


		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")

		users, ok := data["users"].([]interface{})
		assert.True(t, ok, "Response should contain users array")
		assert.Len(t, users, 1, "Users list should contain one user")

		pagination, ok := data["pagination"].(map[string]interface{})
        assert.True(t, ok)
        assert.Equal(t, float64(1), pagination["total_items"])
        assert.Equal(t, float64(1), pagination["current_page"])
	})

	var createdUserID uint

	t.Run("Setup For Get/Update/Delete", func(t *testing.T) {
		testutils.CleanupDatabase(db)
		user := models.User{Name: "Specific User", Email: "specific@example.com"}
		result := db.Create(&user)
		assert.NoError(t, result.Error)
		createdUserID = user.ID
		assert.NotZero(t, createdUserID, "Created user ID should not be zero")
	})


	t.Run("Get User", func(t *testing.T) {
		if createdUserID == 0 {
			t.Skip("Skipping Get User as user creation failed in setup")
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/users/%d", createdUserID), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

        responseData, ok := response.Data.(map[string]interface{})
        assert.True(t, ok)
        assert.Equal(t, "Specific User", responseData["name"])
        assert.Equal(t, "specific@example.com", responseData["email"])
	})

	t.Run("Update User", func(t *testing.T) {
		if createdUserID == 0 {
			t.Skip("Skipping Update User as user creation failed in setup")
		}
		updateData := map[string]interface{}{
			"name":  "Updated Specific User",
			"email": "updatedspecific@example.com",
		}
		jsonData, _ := json.Marshal(updateData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/users/%d", createdUserID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

        responseData, ok := response.Data.(map[string]interface{})
        assert.True(t, ok)
        assert.Equal(t, updateData["name"], responseData["name"])
        assert.Equal(t, updateData["email"], responseData["email"])

		var updatedUser models.User
		db.First(&updatedUser, createdUserID)
		assert.Equal(t, updateData["name"], updatedUser.Name)
		assert.Equal(t, updateData["email"], updatedUser.Email)
	})

	t.Run("Delete User", func(t *testing.T) {
		if createdUserID == 0 {
			t.Skip("Skipping Delete User as user creation failed in setup")
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/users/%d", createdUserID), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Nil(t, response.Data, "Response data should be nil for delete")


		var deletedUser models.User
		err = db.First(&deletedUser, createdUserID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

}

func TestCreateUser_Conflict(t *testing.T) {
	r, userService, db := testutils.SetupTestRouter(true)
	if db == nil { t.Fatal("DB nil") }
	if userService == nil { t.Fatal("userService nil") }
	testutils.CleanupDatabase(db)
	defer testutils.CleanupDatabase(db)

	userHandler := handlers.NewUserHandler(userService)
	r.POST("/api/users", userHandler.Create)

	initialUserData := map[string]interface{}{"name": "Initial User", "email": "conflict@example.com"}
	jsonData, _ := json.Marshal(initialUserData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	conflictUserData := map[string]interface{}{"name": "Conflict User", "email": "conflict@example.com"}
	jsonData, _ = json.Marshal(conflictUserData)
	wConflict := httptest.NewRecorder()
	reqConflict, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	reqConflict.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wConflict, reqConflict)

	assert.Equal(t, http.StatusConflict, wConflict.Code)

	var response utils.Response
	err := json.Unmarshal(wConflict.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, services.ErrUserEmailExists.Error(), response.Message)
}

func TestGetUser_NotFound(t *testing.T) {
    r, userService, db := testutils.SetupTestRouter(true)
	if db == nil { t.Fatal("DB nil") }
	if userService == nil { t.Fatal("userService nil") }
    testutils.CleanupDatabase(db)
    defer testutils.CleanupDatabase(db)

    userHandler := handlers.NewUserHandler(userService)
    r.GET("/api/users/:id", userHandler.Get) // Only need GET for this test case

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/users/99999", nil) // Non-existent ID
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)

    var response utils.Response
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.False(t, response.Success)
    assert.Equal(t, services.ErrUserNotFound.Error(), response.Message) // Check specific error message
}
