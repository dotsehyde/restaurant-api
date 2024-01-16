package server

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetFoods(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get All Food",
	})
}

func (s *Server) GetFood(c *gin.Context) {}

func (s *Server) CreateFood(c *gin.Context) {}

func (s *Server) UpdateFood(c *gin.Context) {}

func (s *Server) DeleteFood(c *gin.Context) {}

func round(num float64) int {
	return int(num + 0.5)
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
