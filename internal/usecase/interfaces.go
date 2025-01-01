package usecase

import (
	"context"

	"github.com/signalable/qauth/internal/domain"
)

// AuthUseCase 인터페이스 정의
type AuthUseCase interface {
	// 토큰 생성
	CreateToken(ctx context.Context, userID string) (*domain.AuthResponse, error)

	// 토큰 검증
	ValidateToken(ctx context.Context, token string) (*domain.TokenValidationResponse, error)

	// 토큰 폐기 (로그아웃)
	RevokeToken(ctx context.Context, token string) error

	// 사용자의 모든 토큰 폐기 (전체 로그아웃)
	RevokeAllTokens(ctx context.Context, userID string) error

	// 토큰 새로고침
	RefreshToken(ctx context.Context, oldToken string) (*domain.AuthResponse, error)

	// 인증 토큰의 메타데이터 조회
	GetTokenMetadata(ctx context.Context, token string) (*domain.TokenMetadata, error)
}
