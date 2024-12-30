package handler

import (
	"encoding/json"
	"net/http"

	"github.com/signalable/qauth/internal/domain"
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

// Register 회원가입 핸들러
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 형식입니다", http.StatusBadRequest)
		return
	}

	if err := h.authUseCase.Register(r.Context(), &req); err != nil {
		switch err {
		case domain.ErrEmailAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "내부 서버 오류", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "회원가입이 완료되었습니다",
	})
}

// Login 로그인 핸들러
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 형식입니다", http.StatusBadRequest)
		return
	}

	resp, err := h.authUseCase.Login(r.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "내부 서버 오류", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Logout 로그아웃 핸들러
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "토큰이 없습니다", http.StatusUnauthorized)
		return
	}

	// "Bearer " 접두사 제거
	tokenString := token[7:]

	if err := h.authUseCase.RevokeToken(r.Context(), tokenString); err != nil {
		http.Error(w, "로그아웃 처리 중 오류가 발생했습니다", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "로그아웃되었습니다",
	})
}
