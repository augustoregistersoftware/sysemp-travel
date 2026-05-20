package main

import (
	"os"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "redis:6379"),
	})

	defer redisClient.Close()

	rateLimiter := middleware.NewRateLimiter(
		redisClient,
		10,
		time.Minute,
	)

	server.Use(middleware.RateLimiterMiddleware(rateLimiter))

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
