package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetAllFood(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get All Food",
	})
}
