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

var tableCollection = "table"

func (s *Server) GetTables(c *gin.Context) {
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
	cur, err := s.db.GetDB().Collection(tableCollection).Find(ctx, bson.M{},
		options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	count, err := s.db.GetDB().Collection(tableCollection).CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	var data []models.Table
	if err := cur.All(ctx, &data); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"totalPage": int64(math.Ceil(float64(count) / float64(limit))),
		"hasMore":   count > int64(limit*page),
		"data":      data,
		"page":      page,
	})

}

func (s *Server) GetTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")
	var table models.Table
	if err := s.db.GetDB().Collection(tableCollection).
		FindOne(ctx, bson.M{"tableId": id}).Decode(&table); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, table)
}

func (s *Server) CreateTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)

	defer cancel()
	var table models.Table
	if err := c.BindJSON(&table); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(&table); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.ID = primitive.NewObjectID()
	table.TableID = table.ID.Hex()

	res, err := s.db.GetDB().Collection(tableCollection).InsertOne(ctx, &table)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (s *Server) UpdateTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)

	defer cancel()

	var table models.Table
	if err := c.BindJSON(&table); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	id := c.Param("id")
	var updateObj primitive.D
	if table.NumberOfGuests != nil {
		updateObj = append(updateObj, bson.E{Key: "numberOfGuests", Value: table.NumberOfGuests})
	}
	if table.TableNumber != nil {
		updateObj = append(updateObj, bson.E{Key: "tableNumber", Value: table.TableNumber})
	}
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: table.UpdatedAt})
	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	res, err := s.db.GetDB().Collection(tableCollection).
		UpdateOne(ctx, bson.M{"tableId": id}, bson.M{"$set": updateObj}, &opt)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (s *Server) DeleteTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")
	res, err := s.db.GetDB().Collection(tableCollection).DeleteOne(ctx, bson.M{"tableId": id})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if res.DeletedCount <= 0 {
		c.JSON(404, gin.H{"message": "Document not found"})
		return
	}
	c.JSON(200, res)
}
