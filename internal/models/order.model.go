package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	OrderDate time.Time          `bson:"orderDate" json:"orderDate" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	OrderID   string             `bson:"orderId" json:"orderId" validate:"required"`
	TableID   string             `bson:"tableId" json:"tableId"`
}

func (o *Order) UpdateUpdatedAt() {
	o.UpdatedAt = time.Now()
}
