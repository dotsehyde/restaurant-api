package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndex(db *mongo.Database) {
	tableNumberIndex := mongo.IndexModel{
		Keys:    bson.M{"tableNumber": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := db.Collection("table").Indexes().CreateOne(context.TODO(), tableNumberIndex)
	if err != nil {
		log.Fatal(err)
	}
}
