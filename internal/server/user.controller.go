package server

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var userCollection = "users"

func (s *Server) GetUsers(c *gin.Context) {
	records, err := s.db.GetDB().Collection(userCollection).Find(c, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"users": records,
	})
}

func (s *Server) GetUser(c *gin.Context) {

}

func (s *Server) LoginUser(c *gin.Context) {}

func (s *Server) SignupUser(c *gin.Context) {}

func hashPassword(password string) string {
	return ""
}

func verifyPassword(hashedPassword, password string) bool {
	return false
}
