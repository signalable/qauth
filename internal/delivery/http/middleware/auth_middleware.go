package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/signalable/qauth/internal/usecase"
)

type AuthMiddleware struct {
	authUseCase usecase.AuthUseCase
}

// NewAuthMiddleware Auth 미들웨어 생성자
func NewAuthMiddleware(authUseCase usecase.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

// Authenticate 인증 미들웨어
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
			return
		}

		// Bearer 토큰 검증
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "잘못된 인증 형식입니다", http.StatusUnauthorized)
			return
		}

		token, err := m.authUseCase.ValidateToken(r.Context(), tokenParts[1])
		if err != nil {
			http.Error(w, "유효하지 않은 토큰입니다", http.StatusUnauthorized)
			return
		}

		// Context에 사용자 ID 추가
		ctx := context.WithValue(r.Context(), "user_id", token.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
