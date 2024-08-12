package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(uri string) *mongo.Client {
	if uri == "" {
		log.Fatal("Set your 'db_uri' on config.json")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// indexModel := mongo.IndexModel{
	// 	Keys:    bson.D{{"sysname", 1}},
	// 	Options: options.Index().SetUnique(true),
	// }

	// if _, err = client.Database(constants.DATABASE).
	// 	Collection(constants.DEVICES_COLLECTION).
	// 	Indexes().
	// 	CreateOne(context.TODO(), indexModel); err != nil {
	// 	panic(err)
	// }

	return client

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
}
