package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Service struct {
	secretKey []byte
}

// NewJWTService JWT 서비스 생성자
func NewJWTService(secretKey string) *Service {
	return &Service{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken JWT 토큰 생성
func (s *Service) GenerateToken(userID string) (string, error) {
	// 토큰 claims 설정
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 토큰 서명
	return token.SignedString(s.secretKey)
}

// ValidateToken JWT 토큰 검증
func (s *Service) ValidateToken(tokenString string) (string, error) {
	// 토큰 파싱 및 검증
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 알고리즘 검증 추가
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("token parsing error: %w", err)
	}

	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	// Claims에서 user_id 추출
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrInvalidKeyType
	}

	// 만료 시간 검증 추가
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return "", fmt.Errorf("token is expired")
		}
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", jwt.ErrInvalidKeyType
	}

	return userID, nil
}
