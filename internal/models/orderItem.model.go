package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Quantity    string             `bson:"quantity" json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice   float64            `bson:"unitPrice" json:"unitPrice"`
	FoodID      string             `bson:"foodId" json:"foodId" validate:"required"`
	OrderItemID string             `bson:"orderItemId" json:"orderItemId"`
	OrderID     string             `bson:"orderId" json:"orderId" validate:"required"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (o *OrderItem) UpdateUpdatedAt() {
	o.UpdatedAt = time.Now()
}
