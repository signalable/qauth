package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/signalable/qauth/internal/domain"
)

type tokenRepository struct {
	client *redis.Client
}

// NewTokenRepository Redis 토큰 레포지토리 생성자
func NewTokenRepository(client *redis.Client) *tokenRepository {
	return &tokenRepository{
		client: client,
	}
}

// Store 토큰 저장
func (r *tokenRepository) Store(ctx context.Context, token *domain.Token) error {
	// 토큰을 JSON으로 직렬화
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}

	// Redis에 토큰 저장 (만료시간 설정)
	duration := time.Until(token.ExpiresAt)
	return r.client.Set(ctx, token.TokenString, tokenJSON, duration).Err()
}

// Verify 토큰 검증
func (r *tokenRepository) Verify(ctx context.Context, tokenString string) (*domain.Token, error) {
	// Redis에서 토큰 조회
	tokenJSON, err := r.client.Get(ctx, tokenString).Result()
	if err == redis.Nil {
		return nil, domain.ErrInvalidToken
	}
	if err != nil {
		return nil, err
	}

	// JSON을 토큰 객체로 역직렬화
	var token domain.Token
	if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
		return nil, err
	}

	// 토큰 유효성 검사
	if token.IsRevoked {
		return nil, domain.ErrRevokedToken
	}
	if time.Now().After(token.ExpiresAt) {
		return nil, domain.ErrExpiredToken
	}

	return &token, nil
}

// Revoke 토큰 폐기
func (r *tokenRepository) Revoke(ctx context.Context, tokenString string) error {
	// 토큰 조회
	token, err := r.Verify(ctx, tokenString)
	if err != nil {
		return err
	}

	// 토큰 폐기 상태로 변경
	token.IsRevoked = true
	return r.Store(ctx, token)
}
