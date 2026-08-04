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
	router.POST("/users/", handler.Register)
	userGroup.Use(middleware.Auth(validator))
	userGroup.GET("/me", handler.GetMe)
	userGroup.PUT("/me", handler.UpdateUser)
	userGroup.DELETE("/me", handler.DeleteUser)

	authGroup := router.Group("/auth")
	authGroup.Use(middleware.Auth(validator))
	router.POST("/auth/login", handler.Login)
	authGroup.POST("/refresh", handler.Refresh)
	authGroup.POST("/logout", handler.Logout)
	authGroup.POST("/logout-all", handler.LogoutAll)

	accountGroup := router.Group("/accounts")
	accountGroup.Use(middleware.Auth(validator))
	accountGroup.POST("", handler.AccountCreate)
	accountGroup.GET("/:accountId", handler.AccountGet)
	accountGroup.GET("/:accountId/statement", handler.GetStatement)
	accountGroup.POST("/:accountId/deposit", handler.AccountDeposit)
	accountGroup.POST("/:accountId/withdraw", handler.AccountWithdraw)
	accountGroup.POST("/:accountId/check-balance", handler.AccountCheckBalance)

	assetGroup := router.Group("/assets")
	assetGroup.GET("", handler.AssetGetAll)
	assetGroup.GET("/:ticker", handler.AssetGetByTicker)
	assetGroup.GET("/:ticker/prices", handler.AssetGetPricesByTicker)
	assetGroup.Use(middleware.Auth(validator))
	assetGroup.POST("", handler.AssetCreateAsset)
	assetGroup.PUT("/:ticker/price", handler.AssetUpdateAssetPriceByTicker)
	assetGroup.PUT("/:id", handler.AssetUpdateAsset)
	assetGroup.DELETE("/:id", handler.AssetDelete)

}
