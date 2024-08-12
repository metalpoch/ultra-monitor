package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/constants"
	"github.com/metalpoch/olt-blueprint/update/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type trafficRepository struct {
	client *mongo.Client
}

func NewTrafficRepository(client *mongo.Client) *trafficRepository {
	return &trafficRepository{client}
}

func (repo trafficRepository) Create(ctx context.Context, traffic *entity.TrafficOLT) (string, error) {
	col := repo.client.Database(constants.DATABASE).Collection(constants.TRAFFIC_OLT_COLLECTION)
	res, err := col.InsertOne(ctx, traffic)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, err
}
