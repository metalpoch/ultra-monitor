package repository

import (
	"context"
	"fmt"

	"github.com/metalpoch/olt-blueprint/auth/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type exampleRepository struct {
	client *mongo.Client
}

func NewExampleRepository(client *mongo.Client) *exampleRepository {
	return &exampleRepository{client}
}

func (repo exampleRepository) Create(ctx context.Context, data *entity.Example) (string, error) {
	//Supongamos que nos conectamos a la DB y almacenamos la entidad data
	return "guardado...", nil
}

func (repo exampleRepository) Get(ctx context.Context, id uint8) (*entity.ExampleResponse, error) {
	//Supongamos que nos conectamos a la DB y buscamos la el id
	message := fmt.Sprintf("Keloke, el ID que buscas es %d", id)
	return &entity.ExampleResponse{Message: message}, nil
}
