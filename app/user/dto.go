package user

import (
	"strings"
	"time"
)

type RegistrationBody struct {
	Name     string `json:"name"     validate:"required,min=2,max=24,alpha_space"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (b *RegistrationBody) Normalize() {
	b.Name = strings.TrimSpace(b.Name)
	b.Email = strings.TrimSpace(strings.ToLower(b.Email))
	b.Password = strings.TrimSpace(b.Password)
}

type LoginBody struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (b *LoginBody) Normalize() {
	b.Email = strings.TrimSpace(strings.ToLower(b.Email))
	b.Password = strings.TrimSpace(b.Password)
}

type AuthResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAuthResponse(user *User) AuthResponse {
	return AuthResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
