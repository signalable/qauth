package usecase

import (
	"context"
	"time"

	"github.com/signalable/qauth/internal/domain"
	"github.com/signalable/qauth/internal/repository"
	"github.com/signalable/qauth/pkg/hash"
	"github.com/signalable/qauth/pkg/jwt"
)

type authUseCase struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtService *jwt.Service
}

// NewAuthUseCase Auth 유스케이스 생성자
func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtService *jwt.Service,
) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
	}
}

// Register 회원가입 구현
func (uc *authUseCase) Register(ctx context.Context, req *domain.RegisterRequest) error {
	// 이메일 중복 체크
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrEmailAlreadyExists
	}

	// 비밀번호 해싱
	hashedPassword, err := hash.GenerateHash(req.Password)
	if err != nil {
		return err
	}

	// 사용자 생성
	user := &domain.User{
		Email:      req.Email,
		Password:   hashedPassword,
		Name:       req.Name,
		IsVerified: false,  // 이메일 인증 전
		Role:       "user", // 기본 역할
	}

	return uc.userRepo.Create(ctx, user)
}

// Login 로그인 구현
func (uc *authUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// 사용자 조회
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// 비밀번호 검증
	if !hash.CompareHash(user.Password, req.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// JWT 토큰 생성
	tokenString, err := uc.jwtService.GenerateToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Redis에 토큰 저장
	token := &domain.Token{
		UserID:      user.ID.Hex(),
		TokenString: tokenString,
		ExpiresAt:   time.Now().Add(24 * time.Hour), // 24시간 유효
		IsRevoked:   false,
	}

	if err := uc.tokenRepo.Store(ctx, token); err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken: tokenString,
		User:        *user,
	}, nil
}

// ValidateToken 토큰 검증 구현
func (uc *authUseCase) ValidateToken(ctx context.Context, tokenString string) (*domain.Token, error) {
	// Redis에서 토큰 검증
	token, err := uc.tokenRepo.Verify(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	// JWT 토큰 검증
	if _, err := uc.jwtService.ValidateToken(tokenString); err != nil {
		return nil, domain.ErrInvalidToken
	}

	return token, nil
}

// RevokeToken 토큰 폐기 구현 (로그아웃)
func (uc *authUseCase) RevokeToken(ctx context.Context, tokenString string) error {
	return uc.tokenRepo.Revoke(ctx, tokenString)
}
