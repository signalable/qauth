package domain

// AuthRequest 인증 요청 DTO
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse 인증 응답 DTO
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// TokenValidationResponse 토큰 검증 응답
type TokenValidationResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
}

// TokenMetadata 토큰 메타데이터
type TokenMetadata struct {
	UserID    string
	IssuedAt  int64
	ExpiresAt int64
}
