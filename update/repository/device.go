package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/constants"
	"github.com/metalpoch/olt-blueprint/update/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type deviceRepository struct {
	client *mongo.Client
}

func NewDeviceRepository(client *mongo.Client) *deviceRepository {
	return &deviceRepository{client}
}

func (repo deviceRepository) Create(ctx context.Context, device *entity.Device) (string, error) {
	col := repo.client.Database(constants.DATABASE).Collection(constants.DEVICES_COLLECTION)
	cursor, err := col.InsertOne(ctx, device)
	if err != nil {
		return "", err
	}
	id := cursor.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (repo deviceRepository) FindAllOffset(ctx context.Context, limit int64, offset int64) ([]*entity.Device, error) {
	devices := []*entity.Device{}
	col := repo.client.Database(constants.DATABASE).Collection(constants.DEVICES_COLLECTION)
	cursor, err := col.Find(ctx, bson.M{}, options.Find().SetLimit(limit).SetSkip(offset))
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func (repo deviceRepository) FindAll(ctx context.Context) ([]*entity.Device, error) {
	devices := []*entity.Device{}
	col := repo.client.Database(constants.DATABASE).Collection(constants.DEVICES_COLLECTION)
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}
