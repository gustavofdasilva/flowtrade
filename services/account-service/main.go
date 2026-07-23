package main

import (
	"account-service/internal/adapters/database"
	"account-service/internal/adapters/handlers/http"
	"account-service/internal/adapters/repository"
	"account-service/internal/adapters/security"
	"account-service/internal/config"
	"account-service/internal/core/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Erro no banco de dados: %v", err)
	}
	defer db.Close()

	passwordHasher := security.NewBcryptHasher(12)
	tokenProvider := security.NewJWTTokenProvider([]byte(cfg.Auth.JWTSecret), cfg.Auth.JWTTokenDuration)

	userRepo := repository.NewPostgresUserRepository(db)
	authRepo := repository.NewPostgresAuthRepository(db)

	userService := services.NewUserService(
		userRepo,
		passwordHasher,
	)

	authService := services.NewAuthService(
		userRepo,
		authRepo,
		passwordHasher,
		tokenProvider,
		cfg.Auth.RefreshTokenDuration,
	)

	userHandler := http.NewUserHandler(userService)
	authHandler := http.NewAuthHandler(authService)

	router := gin.Default()
	http.RegisterUserRoutes(router, userHandler)
	http.RegisterAuthRoutes(router, authHandler)

	log.Printf("Server running at %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("error to start HTTP server: %v", err)
	}
}
