package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	FirstName    *string            `bson:"firstName" json:"firstName" validate:"required,min=3,max=100"`
	LastName     *string            `bson:"lastName" json:"lastName" validate:"required,min=3,max=100"`
	Password     *string            `bson:"password" json:"password" validate:"required,min=6"`
	Email        *string            `bson:"email" json:"email" validate:"email,required"`
	Avatar       *string            `bson:"avatar" json:"avatar"`
	Phone        *string            `bson:"phone" json:"phone" validate:"required"`
	Token        *string            `bson:"token" json:"token"`
	RefreshToken *string            `bson:"refreshToken" json:"refreshToken"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	UserID       string             `bson:"userId" json:"userId"`
}

func (u *User) UpdateUpdatedAt() {
	u.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
}
