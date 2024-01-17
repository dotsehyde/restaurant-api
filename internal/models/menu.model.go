package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Category  string             `bson:"category" json:"category" validate:"required"`
	StartDate *time.Time         `bson:"startDate" json:"startDate" validate:"required"`
	EndDate   *time.Time         `bson:"endDate" json:"endDate" validate:"required"`
	MenuID    string             `bson:"menuId" json:"menuId" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (m *Menu) UpdateUpdatedAt() {
	m.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
}
