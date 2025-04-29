// routes/auth.go
package routes

import (
	"github.com/Caknoooo/go-gin-clean-starter/controller"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func Auth(route *gin.Engine, injector *do.Injector) {
	authController := do.MustInvoke[controller.AuthController](injector)

	routes := route.Group("/api/auth")
	{
		routes.POST("/refresh", authController.RefreshToken)
		routes.POST("/logout", authController.Logout)
	}
}
