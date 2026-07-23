package main

import (
	grpcclient "api-gateway/internal/adapters/grpc/client"
	"api-gateway/internal/adapters/http/handler"
	"api-gateway/internal/adapters/http/routes"
	"api-gateway/internal/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// clientes gRPC
	accountClient := grpcclient.NewAccountClient(cfg.AccountServiceAddr)
	// orderClient := grpcclient.NewOrderClient(cfg.OrderServiceAddr)
	// matchingClient := grpcclient.NewMatchingClient(cfg.MatchingServiceAddr)
	// notifClient := grpcclient.NewNotificationClient(cfg.NotifServiceAddr)

	// handlers
	accountHandler := handler.NewAccountHandler(accountClient)
	// orderHandler := handler.NewOrderHandler(orderClient)
	// matchingHandler := handler.NewMatchingHandler(matchingClient)
	// notifHandler := handler.NewNotificationHandler(notifClient)

	// router
	// r.Use(middleware.Logging(), middleware.RateLimit())

	router := gin.Default()
	routes.RegisterAccountRoutes(router, accountHandler)

	log.Printf("Server running at %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("error to start HTTP server: %v", err)
	}
}
