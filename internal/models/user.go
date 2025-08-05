package models

import (
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	Username  string        `json:"username" bson:"username"`
	Email     string        `json:"email" bson:"email"`
	Password  string        `json:"password" bson:"password"`
	CreatedAt bson.DateTime `json:"createdAt" bson:"created_at"`
}

type PublicUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

type GetAllUsersQuery struct {
	Page int `query:"page,default:1" validate:"min=1"`
}

type GetUserByIDParams struct {
	ID string `uri:"id" validate:"required"`
}

type RegisterUserBody struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func (b *RegisterUserBody) Transform() {
	b.Email = strings.ToLower(b.Email)
	b.Username = strings.ToLower(b.Username)
}

type LoginUserBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func (b *LoginUserBody) Transform() {
	b.Email = strings.ToLower(b.Email)
}
