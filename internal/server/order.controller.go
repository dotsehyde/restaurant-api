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

var orderCollection = "order"

func (s *Server) GetOrders(c *gin.Context) {
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
	cur, err := s.db.GetDB().Collection(orderCollection).
		Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	count, err := s.db.GetDB().Collection(orderCollection).CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	var orders []models.Order
	if err := cur.All(context.TODO(), &orders); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"data":      orders,
		"page":      page,
		"hasMore":   count > int64(page*limit),
		"totalPage": int64(math.Ceil(float64(count) / float64(limit))),
	})

}

func (s *Server) GetOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")

	var invoice models.Invoice
	if err := s.db.GetDB().Collection(orderCollection).
		FindOne(ctx, bson.M{"invoiceId": id}).Decode(&invoice); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, invoice)
}

func (s *Server) CreateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(&order); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	//check table
	var table models.Table
	if err := s.db.GetDB().Collection(tableCollection).
		FindOne(ctx, bson.M{"tableId": order.TableID}).Decode(&table); err != nil {
		c.JSON(500, gin.H{"message": "Table ID does not exist"})
		return
	}
	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	res, err := s.db.GetDB().Collection(orderCollection).InsertOne(ctx, &order)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)

}

func (s *Server) UpdateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")

	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	var updateObj primitive.D
	//check table
	if order.TableID != nil {
		var table models.Table
		if err := s.db.GetDB().Collection(tableCollection).FindOne(ctx, bson.M{"tableId": order.TableID}).Decode(&table); err != nil {
			c.JSON(400, gin.H{"message": "Table ID does not exist"})
			return
		}
		updateObj = append(updateObj, bson.E{Key: "tableId", Value: table.TableID})
	}

	if order.OrderDate != nil {
		*order.OrderDate, _ = time.Parse(time.RFC3339, order.OrderDate.Format(time.RFC3339))
	}
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: order.UpdatedAt})
	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	res, err := s.db.GetDB().Collection(orderCollection).UpdateOne(ctx, bson.M{"orderId": id}, bson.M{"$set": updateObj}, &opt)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (s *Server) DeleteOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")

	res, err := s.db.GetDB().Collection(orderCollection).DeleteOne(ctx, bson.M{"orderId": id})
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if res.DeletedCount <= 0 {
		c.JSON(500, gin.H{
			"message": "Document not found",
		})
		return
	}
	c.JSON(200, res)
}
