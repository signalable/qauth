package usecase

import (
	"context"
	"time"

	"github.com/signalable/qauth/internal/domain"
	"github.com/signalable/qauth/internal/repository"
	"github.com/signalable/qauth/pkg/jwt"
)

type authUseCase struct {
	tokenRepo  repository.TokenRepository
	jwtService *jwt.Service
}

// NewAuthUseCase Auth 유스케이스 생성자
func NewAuthUseCase(
	tokenRepo repository.TokenRepository,
	jwtService *jwt.Service,
) AuthUseCase {
	return &authUseCase{
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
	}
}

// CreateToken 토큰 생성
func (uc *authUseCase) CreateToken(ctx context.Context, userID string) (*domain.AuthResponse, error) {
	// JWT 토큰 생성
	tokenString, err := uc.jwtService.GenerateToken(userID)
	if err != nil {
		return nil, err
	}

	// 토큰 메타데이터 생성
	now := time.Now()
	metadata := &domain.TokenMetadata{
		UserID:    userID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(24 * time.Hour).Unix(),
	}

	// Redis에 토큰 메타데이터 저장
	if err := uc.tokenRepo.Store(ctx, userID, metadata); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   metadata.ExpiresAt - metadata.IssuedAt,
	}, nil
}

// ValidateToken 토큰 검증
func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*domain.TokenValidationResponse, error) {
	// JWT 토큰 검증
	_, err := uc.jwtService.ValidateToken(token)
	if err != nil {
		return &domain.TokenValidationResponse{Valid: false}, err
	}

	// Redis에서 토큰 메타데이터 검증
	metadata, err := uc.tokenRepo.Validate(ctx, token)
	if err != nil {
		return &domain.TokenValidationResponse{Valid: false}, err
	}

	return &domain.TokenValidationResponse{
		Valid:  true,
		UserID: metadata.UserID,
	}, nil
}

// RevokeToken 토큰 폐기
func (uc *authUseCase) RevokeToken(ctx context.Context, token string) error {
	return uc.tokenRepo.Revoke(ctx, token)
}

// RevokeAllTokens 사용자의 모든 토큰 폐기
func (uc *authUseCase) RevokeAllTokens(ctx context.Context, userID string) error {
	return uc.tokenRepo.RevokeAll(ctx, userID)
}

// RefreshToken 토큰 새로고침
func (uc *authUseCase) RefreshToken(ctx context.Context, oldToken string) (*domain.AuthResponse, error) {
	// Redis에서 토큰 새로고침
	metadata, err := uc.tokenRepo.Refresh(ctx, oldToken)
	if err != nil {
		return nil, err
	}

	// 새로운 JWT 토큰 생성
	tokenString, err := uc.jwtService.GenerateToken(metadata.UserID)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   metadata.ExpiresAt - metadata.IssuedAt,
	}, nil
}

// GetTokenMetadata 토큰 메타데이터 조회
func (uc *authUseCase) GetTokenMetadata(ctx context.Context, token string) (*domain.TokenMetadata, error) {
	return uc.tokenRepo.Validate(ctx, token)
}
