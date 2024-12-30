package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/signalable/qauth/internal/config"
	"github.com/signalable/qauth/internal/delivery/http/handler"
	"github.com/signalable/qauth/internal/delivery/http/middleware"
	"github.com/signalable/qauth/internal/delivery/http/routes"
	"github.com/signalable/qauth/internal/repository/mongodb"
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

	// MongoDB 연결
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatalf("MongoDB 연결 실패: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	// MongoDB 연결 테스트
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB 연결 테스트 실패: %v", err)
	}

	// Redis 연결
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Redis 연결 테스트
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis 연결 실패: %v", err)
	}

	// 레포지토리 초기화
	userRepo := mongodb.NewUserRepository(mongoClient.Database(cfg.MongoDB.Database))
	tokenRepo := redisRepository.NewTokenRepository(redisClient)

	// JWT 서비스 초기화
	jwtService := jwt.NewJWTService(cfg.JWT.SecretKey)

	// 유스케이스 초기화
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)

	// 핸들러 및 미들웨어 초기화
	authHandler := handler.NewAuthHandler(authUseCase)
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// 라우터 설정
	router := mux.NewRouter()
	routes.SetupAuthRoutes(router, authHandler, authMiddleware)

	// 서버 시작
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("서버 시작: %s", serverAddr)

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
