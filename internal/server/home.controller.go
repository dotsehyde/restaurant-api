package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HelloWorldHandler(c *gin.Context) {
	// s.db.GetDB().Collection("test").InsertOne(c, map[string]string{"name": "test"})
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}
