package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"

	"github.com/signalable/qauth/internal/config"
	"github.com/signalable/qauth/internal/delivery/http/handler"
	"github.com/signalable/qauth/internal/delivery/http/middleware"
	"github.com/signalable/qauth/internal/delivery/http/routes"
	redisRepository "github.com/signalable/qauth/internal/repository/redis"
	"github.com/signalable/qauth/internal/usecase"
	"github.com/signalable/qauth/pkg/jwt"
)

func main() {
	// 설정 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("설정을 로드할 수 없습니다: %v", err)
	}

	// Redis 연결
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Redis 연결 테스트
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis 연결 실패: %v", err)
	}

	// JWT 서비스 초기화
	jwtService := jwt.NewJWTService(cfg.JWT.SecretKey)

	// 레포지토리 초기화
	tokenRepo := redisRepository.NewTokenRepository(redisClient)

	// 유스케이스 초기화
	authUseCase := usecase.NewAuthUseCase(tokenRepo, jwtService)

	// 핸들러 및 미들웨어 초기화
	authHandler := handler.NewAuthHandler(authUseCase)
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// 라우터 설정
	router := mux.NewRouter()
	routes.SetupAuthRoutes(router, authHandler, authMiddleware)

	// CORS 미들웨어 설정
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// 서버 시작
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Auth Service 시작: %s", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
