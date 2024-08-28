package model

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

type NewUser struct {
	Id               string   `json:_id`
	Email            string   `json:"email"`
	Password         string   `json:"password"`
	ChangePassw      bool     `json:"change_password"`
	Fullname         string   `json:"fullname"`
	P00              uint     `json:"p00"`
	IsAdmin          bool     `json:"is_admin"`
	StatesPermission []string `json:"states_permission"`
}

type Users []*NewUser
