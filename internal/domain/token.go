package domain

import "time"

// Token 도메인 모델
type Token struct {
	UserID      string    `json:"user_id"`
	TokenString string    `json:"token"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsRevoked   bool      `json:"is_revoked"`
}

// Token 생성 요청 DTO
type TokenRequest struct {
	UserID   string
	Duration time.Duration
}
