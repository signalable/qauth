package usecase

import (
	"context"

	"github.com/signalable/qauth/internal/domain"
)

// AuthUseCase 인터페이스 정의
type AuthUseCase interface {
	// 회원가입
	Register(ctx context.Context, req *domain.RegisterRequest) error
	// 로그인
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	// 토큰 검증
	ValidateToken(ctx context.Context, tokenString string) (*domain.Token, error)
	// 토큰 폐기 (로그아웃)
	RevokeToken(ctx context.Context, tokenString string) error
}
