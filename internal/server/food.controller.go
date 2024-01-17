package server

import (
	"context"
	"fmt"
	"math"
	"restaurant-api/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var foodCollection string = "food"
var validate = validator.New()

func (s *Server) GetFoods(c *gin.Context) {

}

func (s *Server) GetFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	foodId := c.Param("id")
	var foodModel models.Food
	err := s.db.GetDB().Collection(foodCollection).FindOne(ctx, bson.M{
		"foodId": foodId,
	}).Decode(&foodModel)

	if err != nil {
		c.JSON(500,
			gin.H{
				"message": err.Error(),
			})
		return
	}
	c.JSON(200, foodModel)
}

func (s *Server) CreateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var food models.Food
	var menu models.Menu
	//Bind JSON
	if err := c.BindJSON(&food); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	//Check validation
	if err := validate.Struct(&food); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	//Check if menu exits
	if err := s.db.GetDB().Collection(menuCollection).FindOne(ctx, bson.M{"menuId": food.MenuID}).Decode(&menu); err != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("Menu doesn't exist: %v", err.Error()),
		})
		return
	}
	//Set variables
	menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.ID = primitive.NewObjectID()
	food.FoodID = food.ID.Hex()
	food.Price = toFixed(food.Price, 2)

	res, err := s.db.GetDB().Collection(foodCollection).InsertOne(ctx, &food)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, res)

}

func (s *Server) UpdateFood(c *gin.Context) {}

func (s *Server) DeleteFood(c *gin.Context) {}

func round(num float64) int {
	return int(num + 0.5)
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
