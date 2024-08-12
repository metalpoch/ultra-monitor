package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/constants"
	"github.com/metalpoch/olt-blueprint/update/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type countRepository struct {
	client *mongo.Client
}

func NewCountRepository(client *mongo.Client) *countRepository {
	return &countRepository{client}
}

func (repo countRepository) Create(ctx context.Context, count entity.Count) (string, error) {
	col := repo.client.Database(constants.DATABASE).Collection(constants.TEMPORAL_COUNT_COLLECTION)
	cursor, err := col.InsertOne(ctx, count)
	if err != nil {
		return "", err
	}

	id := cursor.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (repo countRepository) Find(ctx context.Context, olt, interfaceName string) (entity.Count, error) {
	count := new(entity.Count)
	col := repo.client.Database(constants.DATABASE).Collection(constants.TEMPORAL_COUNT_COLLECTION)
	if err := col.FindOne(ctx, bson.M{"olt": olt, "interface": interfaceName}).Decode(count); err != nil {
		return entity.Count{}, err
	}

	return *count, nil
}
