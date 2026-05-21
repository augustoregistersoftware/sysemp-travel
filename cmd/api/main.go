package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sysemp_feed/auth"
	"sysemp_feed/controller"
	"sysemp_feed/db"
	"sysemp_feed/middleware"
	"sysemp_feed/repository"
	"sysemp_feed/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	server := gin.Default()

	// =========================
	// DATABASE
	// =========================
	dbConnection, err := db.ConnectDB()

	if err != nil {
		panic(err)
	}

	defer dbConnection.Close()

	// =========================
	// JWT
	// =========================

	authService := auth.NewService(
		getEnv("JWT_SECRET", "dev-secret-change-me"),
		24*time.Hour,
	)

	// =========================
	// REDIS
	// =========================

	if getBoolEnv("RATE_LIMIT_ENABLED", true) {
		redisClient := redis.NewClient(&redis.Options{
			Addr: getEnv("REDIS_ADDR", "redis:6379"),
		})

		defer redisClient.Close()

		rateLimitRequests := getIntEnv("RATE_LIMIT_REQUESTS", 10)
		rateLimitWindowSeconds := getIntEnv("RATE_LIMIT_WINDOW_SECONDS", 60)

		rateLimiter := middleware.NewRateLimiter(
			redisClient,
			rateLimitRequests,
			time.Duration(rateLimitWindowSeconds)*time.Second,
		)

		fmt.Printf("Rate limiter enabled: %d requests per %d seconds\n", rateLimitRequests, rateLimitWindowSeconds)

		server.Use(middleware.RateLimiterMiddleware(rateLimiter))
	}

	// =========================
	// DEPENDENCY INJECTION
	// =========================
	baseRepository := repository.NewRepository(dbConnection)

	// Repositories
	UserCreateRepository := repository.NewUserRepository(baseRepository)
	userRepository := repository.NewUserRepository(baseRepository)
	// Use Cases
	UserUseCase := usecase.NewUserUseCase(UserCreateRepository)
	authUsecase := usecase.NewAuthUsecase(&userRepository)
	// Controllers
	authController := controller.NewAuthController(
		authService,
		authUsecase,
	)
	UserController := controller.NewUserController(UserUseCase)

	// =========================
	// ROUTES
	// =========================
	server.POST("/login", authController.Login)
	server.POST("/create_user", UserController.CreateUser)
	server.DELETE("/approved_user/:id", UserController.ApproveUser)

	server.Run(":8080")
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return fallback
	}

	return number
}

func getBoolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}
