package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Email            string             `bson:"email"`
	Password         string             `bson:"password"`
	ChangePassw      bool               `bson:"change_password"`
	Fullname         string             `bson:"fullname"`
	P00              uint               `bson:"p00"`
	IsAdmin          bool               `bson:"is_admin"`
	StatesPermission []string           `bson:"states_permission"`
}

type UserResponse struct {
	Fullname string `bson:"fullname"`
	P00      uint   `bason:"p00"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type Users []*User
