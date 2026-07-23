package main

import (
	"account-service/internal/adapters/database"
	"account-service/internal/adapters/grpc/handler"
	"account-service/internal/adapters/repository"
	"account-service/internal/adapters/security"
	"account-service/internal/config"
	"account-service/internal/core/services"
	"fmt"
	"log"
	"net"
	"strings"

	pb "github.com/gustavofdasilva/flowtrade/proto/pb"
	"google.golang.org/grpc"
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

	grpcServer := grpc.NewServer()
	pb.RegisterAccountServiceServer(grpcServer, handler.NewAccountGRPCHandler(authService, userService))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", strings.TrimLeft(cfg.GRPCPort, ":"))) // :50051
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
