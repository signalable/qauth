package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 도메인 모델
type User struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"-" bson:"password"` // JSON 직렬화에서 제외
	Name       string             `json:"name" bson:"name"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	IsVerified bool               `json:"is_verified" bson:"is_verified"`
	Role       string             `json:"role" bson:"role"`
}

// 회원가입 요청 DTO
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

// 로그인 요청 DTO
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// 로그인 응답 DTO
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}
