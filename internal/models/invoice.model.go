package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	OrderID        *string            `bson:"orderId" json:"orderId"`
	InvoiceID      string             `bson:"invoiceId" json:"invoiceId"`
	PaymentMethod  *string            `bson:"paymentMethod" json:"paymentMethod" validate:"required,eq=cash|eq=card" default:"cash"`
	PaymentStatus  *string            `bson:"paymentStatus" json:"paymentStatus" validate:"required,eq=pending|eq=paid" default:"pending"`
	PaymentDueDate *time.Time         `bson:"paymentDueDate" json:"paymentDueDate" validate:"required"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
