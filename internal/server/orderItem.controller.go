package server

import (
	"context"
	"math"
	"restaurant-api/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderItemCollection = "orderItem"

func (s *Server) GetOrderItems(c *gin.Context) {
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
	cur, err := s.db.GetDB().Collection(orderCollection).Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	count, err := s.db.GetDB().Collection(orderItemCollection).CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	var data []models.OrderItem
	if err := cur.All(context.TODO(), &data); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"totalPage": int64(math.Ceil(float64(count) / float64(limit))),
		"page":      page,
		"data":      data,
		"hasMore":   count > int64(page*limit),
	})
	return

}

func (s *Server) GetOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")

	var data models.OrderItem
	if err := s.db.GetDB().Collection(orderItemCollection).FindOne(ctx, bson.M{"orderItemId": id}).Decode(&data); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, data)

}

func (s *Server) CreateOrderItem(c *gin.Context) {}

func (s *Server) UpdateOrderItem(c *gin.Context) {}

func (s *Server) DeleteOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")
	res, err := s.db.GetDB().Collection(orderCollection).DeleteOne(ctx, bson.M{"orderItemId": id})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if res.DeletedCount <= 0 {
		c.JSON(400, gin.H{"message": "Document not found"})
		return
	}
}

func (s *Server) GetOrderItemsByOrder(c *gin.Context) {}
