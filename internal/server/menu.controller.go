package server

import (
	"context"
	"math"
	"restaurant-api/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection = "menu"

func (s *Server) GetMenus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

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
	res, err := s.db.GetDB().Collection(menuCollection).Find(context.TODO(), bson.M{},
		options.Find().SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "name", Value: 1}}))

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	//Count
	count, err := s.db.GetDB().Collection(menuCollection).CountDocuments(ctx, bson.M{})

	//convert results into []bson.M
	var data []models.Menu
	if err := res.All(ctx, &data); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data":      data,
		"totalPage": int64(math.Ceil(float64(count) / float64(limit))),
		"page":      page,
		"hasMore":   count > int64(page*limit),
	})

}

func (s *Server) GetMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var menu models.Menu
	menuId := c.Param("id")
	err := s.db.GetDB().Collection(menuCollection).FindOne(ctx, bson.M{"menuId": menuId}).Decode(&menu)

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, menu)

}

func (s *Server) CreateMenu(c *gin.Context) {
	// fmt.Println(time.Now().Format(time.RFC3339))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var menu models.Menu
	if err := c.BindJSON(&menu); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := validate.Struct(menu); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.ID = primitive.NewObjectID()
	menu.MenuID = menu.ID.Hex()
	res, err := s.db.GetDB().Collection(menuCollection).InsertOne(ctx, &menu)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, res)

}

func (s *Server) UpdateMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var menu models.Menu
	if err := c.BindJSON(&menu); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	menuId := c.Param("id")
	var updateObj primitive.D
	if menu.StartDate != nil || menu.EndDate != nil {
		if !inTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
			c.JSON(400, gin.H{"message": "Invalid date input"})
			return
		}
		*menu.StartDate, _ = time.Parse(time.RFC3339, menu.StartDate.Format(time.RFC3339))
		*menu.EndDate, _ = time.Parse(time.RFC3339, menu.EndDate.Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "startDate", Value: menu.StartDate})
		updateObj = append(updateObj, bson.E{Key: "endDate", Value: menu.EndDate})
	}
	if menu.Name != nil {
		updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
	}
	if menu.Category != nil {
		updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
	}
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: menu.UpdatedAt})
	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	res, err := s.db.GetDB().Collection(menuCollection).UpdateOne(ctx, bson.M{"menuId": menuId}, bson.M{"$set": updateObj}, &opt)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
	return
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}

func (s *Server) DeleteMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	id := c.Param("id")

	res, err := s.db.GetDB().Collection(menuCollection).DeleteOne(ctx, bson.M{"menuId": id})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
}
