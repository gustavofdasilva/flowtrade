package routes

import (
	"api-gateway/internal/adapters/http/handler"
	"api-gateway/internal/adapters/http/middleware"
	"api-gateway/internal/core/ports"

	"github.com/gin-gonic/gin"
)

//TODO: Add auth middleware

func RegisterAccountRoutes(router *gin.Engine, handler *handler.AccountHandler, validator ports.TokenValidator) {

	userGroup := router.Group("/users")
	userGroup.POST("/", handler.Register)

	authGroup := router.Group("/auth")
	authGroup.Use(middleware.Auth(validator))
	router.POST("/auth/login", handler.Login)
	authGroup.POST("/refresh", handler.Refresh)
	authGroup.POST("/logout", handler.Logout)
	authGroup.POST("/logout-all", handler.LogoutAll)

}
