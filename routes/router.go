package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

// RegisterRoutes initializes and registers all application-level routes with the provided Gin engine and dependency injector.
var RegisterRoutes = func(server *gin.Engine, injector *do.Injector) {
	User(server, injector)
}
