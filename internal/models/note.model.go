package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Text      string             `bson:"text" json:"text"`
	Title     string             `bson:"title" json:"title"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	NoteID    string             `bson:"noteId" json:"noteId"`
}

func (n *Note) UpdateUpdatedAt() {
	n.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
}
