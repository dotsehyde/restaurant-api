package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	OrderID        string             `bson:"orderId" json:"orderId" validate:"required"`
	InvoiceID      string             `bson:"invoiceId" json:"invoiceId" validate:"required"`
	PaymentMethod  string             `bson:"paymentMethod" json:"paymentMethod" validate:"required,eq=CASH|eq=CARD" default:"CASH"`
	PaymentStatus  string             `bson:"paymentStatus" json:"paymentStatus" validate:"required,eq=PENDING|eq=PAID" default:"PENDING"`
	PaymentDueDate time.Time          `bson:"paymentDueDate" json:"paymentDueDate" validate:"required"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (i *Invoice) UpdateUpdatedAt() {
	i.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
}
