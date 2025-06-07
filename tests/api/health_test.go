package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/api/handlers"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/pkg/utils"
	"github.com/abhi9s-realm/lean-backend-boilerplate-golang/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Setup - health check doesn't need database
	r, _ := testutils.SetupTestRouter(false)
	r.GET("/api/health", handlers.HealthCheck)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)

	// Perform request
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response utils.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Service is healthy", response.Message)
}
