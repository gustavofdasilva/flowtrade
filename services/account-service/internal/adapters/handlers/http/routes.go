package http

import "github.com/gin-gonic/gin"

//TODO: Add auth middleware

func RegisterUserRoutes(router *gin.Engine, handler *UserHandler) {
	authGroup := router.Group("/users")
	{
		authGroup.POST("/", handler.Register)
		authGroup.GET("/me", handler.GetUserInfo)
		authGroup.PATCH("/me", handler.UpdateUser)
		authGroup.DELETE("/me", handler.DeleteUser)
	}
}

func RegisterAuthRoutes(router *gin.Engine, handler *AuthHandler) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/refresh", handler.RefreshToken)
		authGroup.POST("/logout", handler.Logout)
		authGroup.POST("/logout-all", handler.LogoutAll)
	}
}
