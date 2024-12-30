package repository

import (
	"context"

	"github.com/signalable/qauth/internal/domain"
)

// UserRepository 인터페이스 정의
type UserRepository interface {
	// 사용자 생성
	Create(ctx context.Context, user *domain.User) error
	// 이메일로 사용자 찾기
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	// 사용자 ID로 찾기
	FindByID(ctx context.Context, id string) (*domain.User, error)
	// 이메일 존재 여부 확인
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// 사용자 정보 업데이트
	Update(ctx context.Context, user *domain.User) error
}

// TokenRepository 인터페이스 정의
type TokenRepository interface {
	// 토큰 저장
	Store(ctx context.Context, token *domain.Token) error
	// 토큰 검증
	Verify(ctx context.Context, tokenString string) (*domain.Token, error)
	// 토큰 폐기
	Revoke(ctx context.Context, tokenString string) error
}
