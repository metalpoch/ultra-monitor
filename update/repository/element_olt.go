package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/constants"
	"github.com/metalpoch/olt-blueprint/update/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type oltElementRepository struct {
	client *mongo.Client
}

func NewOltElementRepository(client *mongo.Client) *oltElementRepository {
	return &oltElementRepository{client}
}

func (repo oltElementRepository) Create(ctx context.Context, olt *entity.ElementOLT) (string, error) {
	col := repo.client.Database(constants.DATABASE).Collection(constants.OLT_COLLECTION)
	cursor, err := col.InsertOne(ctx, olt)
	if err != nil {
		return "", err
	}

	id := cursor.InsertedID.(primitive.ObjectID).Hex()
	return id, err
}

func (repo oltElementRepository) FindID(ctx context.Context, olt, interfName string) (string, error) {
	element := new(entity.ElementOLT)
	col := repo.client.Database(constants.DATABASE).Collection(constants.OLT_COLLECTION)
	if err := col.FindOne(ctx, bson.M{"olt": olt, "interface": interfName}).Decode(element); err != nil {
		return "", err
	}

	return element.ID.Hex(), nil
}
