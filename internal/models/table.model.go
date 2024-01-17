package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	NumberOfGuests int                `bson:"numberOfGuests" json:"numberOfGuests" validate:"required"`
	TableNumber    int                `bson:"tableNumber" json:"tableNumber" validate:"required"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	TableID        string             `bson:"tableId" json:"tableId" validate:"required"`
}

func (t *Table) UpdateUpdatedAt() {
	t.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
}
