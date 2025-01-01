package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/signalable/qauth/internal/usecase"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

// NewAuthHandler Auth 핸들러 생성자
func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// CreateToken 토큰 생성 핸들러
func (h *AuthHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID") // User Service에서 전달받은 사용자 ID
	if userID == "" {
		http.Error(w, "사용자 ID가 필요합니다", http.StatusBadRequest)
		return
	}

	resp, err := h.authUseCase.CreateToken(r.Context(), userID)
	if err != nil {
		http.Error(w, "토큰 생성 실패", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ValidateToken 토큰 검증 핸들러
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		http.Error(w, "토큰이 필요합니다", http.StatusBadRequest)
		return
	}

	// 디버깅을 위한 로그 추가
	log.Printf("Received token for validation: %s", token)

	resp, err := h.authUseCase.ValidateToken(r.Context(), token)
	if err != nil {
		log.Printf("Token validation error: %v", err) // 에러 로깅 추가
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RevokeToken 토큰 폐기 핸들러
func (h *AuthHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		http.Error(w, "토큰이 필요합니다", http.StatusBadRequest)
		return
	}

	if err := h.authUseCase.RevokeToken(r.Context(), token); err != nil {
		http.Error(w, "토큰 폐기 실패", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "토큰이 폐기되었습니다",
	})
}

// RefreshToken 토큰 새로고침 핸들러
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		http.Error(w, "토큰이 필요합니다", http.StatusBadRequest)
		return
	}

	resp, err := h.authUseCase.RefreshToken(r.Context(), token)
	if err != nil {
		http.Error(w, "토큰 새로고침 실패", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// extractToken 요청에서 토큰 추출
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
