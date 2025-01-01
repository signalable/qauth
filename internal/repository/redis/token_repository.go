package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/signalable/qauth/internal/domain"
	"github.com/signalable/qauth/pkg/jwt"
)

type tokenRepository struct {
	client     *redis.Client
	jwtService *jwt.Service
}

// NewTokenRepository Redis 토큰 레포지토리 생성자
func NewTokenRepository(client *redis.Client, jwtService *jwt.Service) *tokenRepository {
	return &tokenRepository{
		client:     client,
		jwtService: jwtService,
	}
}

// Store 토큰 저장
func (r *tokenRepository) Store(ctx context.Context, userID string, metadata *domain.TokenMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("토큰 메타데이터 직렬화 실패: %w", err)
	}

	// userID를 키로 사용
	key := fmt.Sprintf("user:%s:token", userID)
	duration := time.Until(time.Unix(metadata.ExpiresAt, 0))

	return r.client.Set(ctx, key, data, duration).Err()
}

// Validate 토큰 검증
func (r *tokenRepository) Validate(ctx context.Context, token string) (*domain.TokenMetadata, error) {
	// Redis에서 토큰 조회 전 로깅 추가
	log.Printf("Validating token in Redis: %s", token)

	// JWT에서 userID를 추출
	userID, err := r.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// 토큰 메타데이터 조회
	data, err := r.client.Get(ctx, fmt.Sprintf("user:%s:token", userID)).Bytes()
	if err == redis.Nil {
		return nil, domain.ErrInvalidToken
	}
	if err != nil {
		return nil, fmt.Errorf("토큰 조회 실패: %w", err)
	}

	var metadata domain.TokenMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("토큰 메타데이터 역직렬화 실패: %w", err)
	}

	// 토큰 만료 검사
	if time.Now().Unix() > metadata.ExpiresAt {
		return nil, domain.ErrExpiredToken
	}

	return &metadata, nil
}

// Revoke 토큰 폐기
func (r *tokenRepository) Revoke(ctx context.Context, token string) error {
	metadata, err := r.Validate(ctx, token)
	if err != nil {
		return err
	}

	// 토큰 삭제
	if err := r.client.Del(ctx, token).Err(); err != nil {
		return fmt.Errorf("토큰 삭제 실패: %w", err)
	}

	// 사용자의 토큰 목록에서 제거
	if err := r.client.SRem(ctx, fmt.Sprintf("user:%s:tokens", metadata.UserID), token).Err(); err != nil {
		return fmt.Errorf("토큰 매핑 삭제 실패: %w", err)
	}

	return nil
}

// RevokeAll 사용자의 모든 토큰 폐기
func (r *tokenRepository) RevokeAll(ctx context.Context, userID string) error {
	// 사용자의 모든 토큰 조회
	tokens, err := r.client.SMembers(ctx, fmt.Sprintf("user:%s:tokens", userID)).Result()
	if err != nil {
		return fmt.Errorf("토큰 목록 조회 실패: %w", err)
	}

	// 각 토큰 삭제
	for _, token := range tokens {
		if err := r.client.Del(ctx, token).Err(); err != nil {
			return fmt.Errorf("토큰 삭제 실패: %w", err)
		}
	}

	// 토큰 목록 삭제
	if err := r.client.Del(ctx, fmt.Sprintf("user:%s:tokens", userID)).Err(); err != nil {
		return fmt.Errorf("토큰 목록 삭제 실패: %w", err)
	}

	return nil
}

// Refresh 토큰 새로고침
func (r *tokenRepository) Refresh(ctx context.Context, oldToken string) (*domain.TokenMetadata, error) {
	metadata, err := r.Validate(ctx, oldToken)
	if err != nil {
		return nil, err
	}

	// 새로운 만료 시간 설정
	metadata.IssuedAt = time.Now().Unix()
	metadata.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	// 새 토큰 저장
	if err := r.Store(ctx, metadata.UserID, metadata); err != nil {
		return nil, err
	}

	// 이전 토큰 폐기
	if err := r.Revoke(ctx, oldToken); err != nil {
		return nil, err
	}

	return metadata, nil
}
