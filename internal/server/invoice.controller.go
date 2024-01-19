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

var invoiceCollection = "invoice"

func (s *Server) GetInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")

	var invoice models.Invoice

	if err := s.db.GetDB().Collection(invoiceCollection).
		FindOne(ctx, bson.M{"invoiceId": id}).Decode(&invoice); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, invoice)
}

func (s *Server) GetInvoices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 1
	}
	skip := (page - 1) * limit
	cur, err := s.db.GetDB().Collection(invoiceCollection).
		Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	count, err := s.db.GetDB().Collection(invoiceCollection).CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	var data []models.Invoice
	if err := cur.All(context.TODO(), &data); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"data":      data,
		"page":      page,
		"hasMore":   count > int64(page*limit),
		"totalPage": int64(math.Ceil(float64(count) / float64(limit))),
	})
}

func (s *Server) CreateInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	var invoice models.Invoice
	if err := c.BindJSON(&invoice); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if err := validate.Struct(&invoice); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	//check order
	var order models.Order
	if err := s.db.GetDB().Collection(orderCollection).FindOne(ctx, bson.M{"orderId": invoice.OrderID}).Decode(&order); err != nil {
		c.JSON(400, gin.H{"message": "Order ID does not exist."})
		return
	}
	invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.ID = primitive.NewObjectID()
	invoice.InvoiceID = invoice.ID.Hex()

	res, err := s.db.GetDB().Collection(invoiceCollection).InsertOne(ctx, &invoice)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (s *Server) UpdateInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")
	var invoice models.Invoice
	var order models.Order
	if err := c.BindJSON(&invoice); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
	}
	var updateObj primitive.D
	if invoice.OrderID != nil {
		//check orderID
		if err := s.db.GetDB().Collection(orderCollection).
			FindOne(ctx, bson.M{"orderId": invoice.OrderID}).Decode(&order); err != nil {
			c.JSON(400, gin.H{"message": "Order ID does not exist"})
			return
		}
		updateObj = append(updateObj, bson.E{Key: "orderId", Value: order.OrderID})
	}

	if invoice.PaymentMethod != nil {
		updateObj = append(updateObj, bson.E{Key: "paymentMethod", Value: invoice.PaymentMethod})
	}
	if invoice.PaymentStatus != nil {
		updateObj = append(updateObj, bson.E{Key: "paymentStatus", Value: invoice.PaymentStatus})
	}
	if invoice.PaymentDueDate != nil {
		*invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, invoice.PaymentDueDate.Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "paymentDueDate", Value: invoice.PaymentDueDate})
	}
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: invoice.UpdatedAt})
	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	res, err := s.db.GetDB().Collection(invoiceCollection).
		UpdateOne(ctx, bson.M{"invoiceId": id}, bson.M{"$set": updateObj}, &opt)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (s *Server) DeleteInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	id := c.Param("id")
	res, err := s.db.GetDB().Collection(invoiceCollection).DeleteOne(ctx, bson.M{"invoiceId": id})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if res.DeletedCount <= 0 {
		c.JSON(404, gin.H{"message": "Document not found"})
		return
	}
	c.JSON(200, res)
}
