package server

import (
	"context"
	"restaurant-api/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var menuCollection = "menu"

func (s *Server) GetMenus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	res, err := s.db.GetDB().Collection(menuCollection).Find(context.TODO(), bson.M{})

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	//convert results into []bson.M
	var data []bson.M

	if err := res.All(ctx, &data); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, data)

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
	// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// defer cancel()

	// var menu models.Menu
	// if err := c.BindJSON(&menu); err != nil {
	// 	c.JSON(400, gin.H{
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }
	// menuId := c.Param("id")

	// updateObj := primitive.D{}

	// if menu.StartDate != nil && menu.EndDate != nil {
	// 	if !inTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
	// 		c.JSON(400, gin.H{
	// 			"message": "Invalid date input",
	// 		})
	// 		return
	// 	}
	// 	updateObj = append(updateObj, bson.E{Key: "startDate", Value: menu.StartDate})
	// 	updateObj = append(updateObj, bson.E{Key: "endDate", Value: menu.EndDate})

	// 	if menu.Name != "" {
	// 		updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
	// 	}
	// 	if menu.Category != "" {
	// 		updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
	// 	}
	// 	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// 	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: menu.UpdatedAt})
	// 	upsert := true
	// 	opt := options.UpdateOptions{
	// 		Upsert: &upsert,
	// 	}
	// 	res, err := s.db.GetDB().Collection(menuCollection).UpdateOne(ctx, bson.M{"menuId": menuId}, bson.D{
	// 		{Key: "$set", Value: updateObj},
	// 	}, &opt)
	// 	if err != nil {
	// 		c.JSON(500, gin.H{
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}
	// 	c.JSON(
	// 		200, res,
	// 	)
	// }
}

func (s *Server) DeleteMenu(c *gin.Context) {}
