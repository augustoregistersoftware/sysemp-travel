package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sysemp_travel/auth"
	"sysemp_travel/config"
	"sysemp_travel/controller"
	"sysemp_travel/db"
	"sysemp_travel/middleware"
	"sysemp_travel/repository"
	"sysemp_travel/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	println(config.GetLogo())
	// =========================
	// Config Gin
	// =========================
	server := gin.Default()
	err := server.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

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
	paymentsRepository := repository.NewPaymentsRepository(baseRepository)
	// Use Cases
	UserUseCase := usecase.NewUserUseCase(UserCreateRepository)
	authUsecase := usecase.NewAuthUsecase(&userRepository)
	PaymentsUseCase := usecase.NewPaymentsUseCase(paymentsRepository)
	// Controllers
	authController := controller.NewAuthController(
		authService,
		authUsecase,
	)
	UserController := controller.NewUserController(UserUseCase)
	PaymentsController := controller.NewPaymentsController(PaymentsUseCase)

	// =========================
	// ROUTES
	// =========================
	server.POST("/login", authController.Login)
	server.POST("/create_user", UserController.CreateUser)
	server.GET("/users", UserController.Users)
	server.GET("/users_approved_list", UserController.UsersApprovedList)

	server.PATCH("/reproved_user/:id", UserController.ReproveUser)
	server.DELETE("/approved_user/:id", UserController.ApproveUser)

	server.GET("/payments", PaymentsController.Payments)

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
