package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/controller"
	"github.com/Caknoooo/go-gin-clean-starter/middleware"
	"github.com/Caknoooo/go-gin-clean-starter/service"
)

// User registers the routes for user-related operations such as registration, login, and profile management.
var User = func(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userController := do.MustInvoke[controller.UserController](injector)

	routes := route.Group("/api/user")
	{
		routes.POST("", userController.Register)
		routes.GET("", userController.GetAllUser)
		routes.POST("/login", userController.Login)
		routes.POST("/refresh", userController.Refresh)
		routes.DELETE("", middleware.Authenticate(jwtService), userController.Delete)
		routes.PATCH("", middleware.Authenticate(jwtService), userController.Update)
		routes.GET("/me", middleware.Authenticate(jwtService), userController.Me)
		routes.POST("/verify_email", userController.VerifyEmail)
		routes.POST("/send_verification_email", userController.SendVerificationEmail)
	}
}
