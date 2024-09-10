package repository

import (
	"context"
	"strconv"

	"github.com/metalpoch/olt-blueprint/auth/constants"
	"github.com/metalpoch/olt-blueprint/auth/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) *userRepository {
	return &userRepository{client}
}

func (repo userRepository) Create(ctx context.Context, data *entity.User) (string, error) {
	var err error
	col := repo.client.Database(constants.DATABASE).Collection(constants.USER_COLLECTION)

	_, err = col.InsertOne(ctx, data)

	if err != nil {
		return "Ocurrio un error al insertar", err
	}
	return "Datos Guardados", nil
}

func (repo userRepository) Get(ctx context.Context) (*entity.Users, error) {
	var users entity.Users
	filter := bson.D{}

	col := repo.client.Database(constants.DATABASE).Collection(constants.USER_COLLECTION)
	cursor, err := col.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user entity.User
		err = cursor.Decode(&user)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return &users, nil
}

func (repo userRepository) GetValue(ctx context.Context, clave string, valor string) (*entity.Users, error) {
	var users entity.Users
	col := repo.client.Database(constants.DATABASE).Collection(constants.USER_COLLECTION)

	filter := bson.M{
		clave: valor,
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user entity.User
		err = cursor.Decode(&user)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return &users, nil
}

func (repo userRepository) DeleteName(ctx context.Context, name string) (string, error) {
	//var err error
	filter := bson.M{"p00": name}
	col := repo.client.Database(constants.DATABASE).Collection(constants.USER_COLLECTION)

	res, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return "Ocurrio un error al intentar eliminar", err
	}
	temp := res.DeletedCount
	// temp := strconv.Itoa(int(res.DeletedCount))

	if temp > 0 {
		return "Se elimino con exito", nil
	} else {
		return strconv.FormatInt(temp, 10), nil
	}

}

func (repo userRepository) ChangePassword(ctx context.Context, user *entity.User) (string, error) {
	//var err error
	col := repo.client.Database(constants.DATABASE).Collection(constants.USER_COLLECTION)
	filter := bson.M{
		"p00":             user.P00,
		"change_password": true,
	}
	update := bson.M{
		"$set": bson.M{
			"password":        user.Password,
			"change_password": false,
		},
	}
	res, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		return "Ocurrio un error al intenar actualizar la contraseña", err
	}

	temp1 := strconv.Itoa(int(res.MatchedCount))
	temp2 := strconv.Itoa(int(res.ModifiedCount))
	temp3 := strconv.Itoa(int(res.UpsertedCount))
	temp4 := temp1 + "," + temp2 + "," + temp3

	if temp4 != "0,0,0" {
		return "Se actualizo la  contraseña correctamente", nil
	} else {
		return temp4, nil
	}
}

func (repo userRepository) Login(ctx context.Context, email string, pasword string) (*entity.User, error) {
	col := repo.client.Database(constants.DATABASE).Collection(constants.USER_COLLECTION)
	var user entity.User
	filter := bson.M{
		"email":    email,
		"password": pasword,
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		err = cursor.Decode(&user)

		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}
