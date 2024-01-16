package tests

import (
	"net/http"
	"net/http/httptest"
	"restaurant-api/internal/server"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllFoods(t *testing.T) {
	s := &server.Server{}
	r := gin.New()
	r.GET("/food/all", s.GetFoods)
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/food/all", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	// Check results
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "{\"message\":\"Get All Food\"}")
}
