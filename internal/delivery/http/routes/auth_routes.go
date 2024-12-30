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
	// 공개 라우트
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// 인증이 필요한 라우트
	router.HandleFunc("/api/auth/logout", authMiddleware.Authenticate(authHandler.Logout)).Methods("POST")
}
