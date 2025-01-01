package routes

import (
	"github.com/gorilla/mux"
	"github.com/signalable/qauth/internal/delivery/http/handler"
	"github.com/signalable/qauth/internal/delivery/http/middleware"
)

// SetupAuthRoutes 라우터 설정
func SetupAuthRoutes(
	router *mux.Router,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// 내부 서비스 간 API (User Service에서 호출)
	router.HandleFunc("/api/auth/token", authHandler.CreateToken).Methods("POST")
	router.HandleFunc("/api/auth/token/validate", authHandler.ValidateToken).Methods("GET")

	// 클라이언트 API
	router.HandleFunc("/api/auth/token/refresh", authHandler.RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/token/revoke", authMiddleware.Authenticate(authHandler.RevokeToken)).Methods("POST")
}
