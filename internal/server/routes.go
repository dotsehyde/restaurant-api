package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	//Home Route '/'
	homeGroup := r.Group("/")
	homeGroup.GET("/", s.HelloWorldHandler)
	homeGroup.GET("/health", s.healthHandler)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
