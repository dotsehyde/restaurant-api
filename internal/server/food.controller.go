package server

import (
	"context"
	"fmt"
	"math"
	"restaurant-api/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection string = "food"
var validate = validator.New()

func (s *Server) GetFoods(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	skip := (page - 1) * limit
	var foods []models.Food

	cur, err := s.db.GetDB().Collection(foodCollection).Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	count, err := s.db.GetDB().Collection(foodCollection).CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	defer cur.Close(ctx)
	if err := cur.All(ctx, &foods); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"page":      page,
		"totalPage": int64(math.Ceil(float64(count) / float64(limit))),
		"hasMore":   count > int64(page*limit),
		"data":      foods})
}

func (s *Server) GetFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)

	foodId := c.Param("id")
	var foodModel models.Food
	err := s.db.GetDB().Collection(foodCollection).FindOne(ctx, bson.M{
		"foodId": foodId,
	}).Decode(&foodModel)
	defer cancel()
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
	food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.ID = primitive.NewObjectID()
	food.FoodID = food.ID.Hex()
	*food.Price = toFixed(*food.Price, 2)

	res, err := s.db.GetDB().Collection(foodCollection).InsertOne(ctx, &food)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, res)
	return
}

func (s *Server) UpdateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	var food models.Food
	var menu models.Menu

	foodId := c.Param("id")
	if err := c.BindJSON(&food); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var updateObj primitive.D
	if food.MenuID != nil {
		// check menu exist
		if err := s.db.GetDB().Collection(menuCollection).FindOne(ctx, bson.M{"menuId": food.MenuID}).Decode(&menu); err != nil {
			c.JSON(400, gin.H{"message": "Menu does not exist"})
			return
		}
		food.MenuID = &menu.MenuID
		updateObj = append(updateObj, bson.E{Key: "menuId", Value: food.MenuID})
	}
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: food.UpdatedAt})
	if food.Name != nil {
		updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
	}
	if food.Price != nil {
		updateObj = append(updateObj, bson.E{Key: "price", Value: toFixed(*food.Price, 2)})
	}
	if food.FoodImage != nil {
		updateObj = append(updateObj, bson.E{Key: "foodImage", Value: food.FoodImage})
	}
	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	res, err := s.db.GetDB().Collection(foodCollection).
		UpdateOne(ctx, bson.M{"foodId": foodId}, bson.M{"$set": updateObj}, &opt)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (s *Server) DeleteFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	id := c.Param("id")

	res, err := s.db.GetDB().Collection(foodCollection).DeleteOne(ctx, bson.M{"foodId": id})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if res.DeletedCount <= 0 {
		c.JSON(404, gin.H{"message": "Document not found"})
		return
	}
	c.JSON(200, res)
	return

}

func round(num float64) int {
	return int(num + 0.5)
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
