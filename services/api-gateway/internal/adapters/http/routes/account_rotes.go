package routes

import (
	"api-gateway/internal/adapters/http/handler"

	"github.com/gin-gonic/gin"
)

//TODO: Add auth middleware

func RegisterAccountRoutes(router *gin.Engine, handler *handler.AccountHandler) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", handler.Register)
	}
}
