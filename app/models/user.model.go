package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name            string     `gorm:"type:varchar(225);not null"`
	Email           string     `gorm:"type:varchar(225);uniqueIndex;not null"`
	Password        string     `gorm:"type:varchar(225);not null"`
	EmailVerifiedAt *time.Time `gorm:"default:null"`
	CreatedAt       *time.Time `gorm:"not null;default:now()"`
	UpdatedAt       *time.Time `gorm:"default:null"`
}

type UserResponse struct {
	ID              *uuid.UUID `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type UserRegisterRequest struct {
	Name            string `json:"name" validate:"required,min=4,max=225"`
	Email           string `json:"email" validate:"required,email,max=225"`
	Password        string `json:"password" validate:"required,min=8,max=30"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8,max=30,eqfield=Password"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func UserFilterRecord(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
