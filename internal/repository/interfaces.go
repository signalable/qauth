package repository

import (
	"context"

	"github.com/signalable/qauth/internal/domain"
)

// TokenRepository 인터페이스 정의
type TokenRepository interface {
	// 토큰 저장
	Store(ctx context.Context, userID string, metadata *domain.TokenMetadata) error

	// 토큰 검증
	Validate(ctx context.Context, token string) (*domain.TokenMetadata, error)

	// 토큰 폐기
	Revoke(ctx context.Context, token string) error

	// 사용자의 모든 토큰 폐기
	RevokeAll(ctx context.Context, userID string) error

	// 토큰 새로고침
	Refresh(ctx context.Context, oldToken string) (*domain.TokenMetadata, error)
}
