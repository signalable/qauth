package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	LogLevel string
}

type ServerConfig struct {
	Host string
	Port string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	SecretKey        string
	ExpirationPeriod time.Duration
}

// LoadConfig .env 파일에서 설정을 로드
func LoadConfig() (*Config, error) {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// Redis DB 번호 파싱
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	// JWT 만료 시간 파싱
	jwtExpirationHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		jwtExpirationHours = 24
	}

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "auth_db"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			SecretKey:        getEnv("JWT_SECRET_KEY", "your-secret-key"),
			ExpirationPeriod: time.Duration(jwtExpirationHours) * time.Hour,
		},
		LogLevel: getEnv("LOG_LEVEL", "debug"),
	}, nil
}

// getEnv 환경 변수를 가져오거나 기본값 반환
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
