package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Price     float64            `bson:"price" json:"price" validate:"required"`
	FoodImage string             `bson:"foodImage" json:"foodImage" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	MenuID    string             `bson:"menuId" json:"menuId" validate:"required"`
	FoodID    string             `bson:"foodId" json:"foodId" validate:"required"`
}

func (f *Food) UpdateUpdatedAt() {
	f.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
}
