package router

import (
	"be-golang-poliklinik/internal/auth"
	"be-golang-poliklinik/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(api *gin.RouterGroup, h *Handlers, authService auth.Service) {
	authRoute := api.Group("/auth")
	{
		authRoute.POST("/register", h.Auth.Register)
		authRoute.POST("/login", h.Auth.Login)
		authRoute.POST("/logout", middleware.AuthMiddleware(authService), h.Auth.Logout)
	}
}
