package models

import "strings"

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
}

type PublicUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

type RegisterUserBody struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

func (b *RegisterUserBody) Normalize() {
	b.Email = strings.ToLower(b.Email)
	b.Username = strings.ToLower(b.Username)
}

type LoginUserBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

func (b *LoginUserBody) Normalize() {
	b.Email = strings.ToLower(b.Email)
}
