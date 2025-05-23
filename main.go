// @title User Management API
// @version 1.0
// @description This is a user management API built with Go, Gin, GORM, and PostgreSQL.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8888
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authentication
package main

import (
	"log"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/Caknoooo/go-gin-clean-starter/command"
	_ "github.com/Caknoooo/go-gin-clean-starter/docs"
	"github.com/Caknoooo/go-gin-clean-starter/middleware"
	"github.com/Caknoooo/go-gin-clean-starter/provider"
	"github.com/Caknoooo/go-gin-clean-starter/routes"
)

// args is responsible for processing command-line arguments and determining whether the application should proceed or exit.
var args = func(injector *do.Injector) bool {
	if len(os.Args) > 1 {
		flag := command.Commands(injector)
		return flag
	}

	return true
}

// run is a variable that defines a function to configure and run a Gin server with the specified routes and settings.
var run = func(server *gin.Engine) {
	server.Static("/assets", "./assets")

	if os.Getenv("IS_LOGGER") == "true" {
		routes.LoggerRoute(server)
	}

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "0.0.0.0:" + port
	} else if os.Getenv("APP_ENV") == "prod" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	myFigure := figure.NewColorFigure("Caknoo", "", "green", true)
	myFigure.Print()

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

// main initializes the application, sets up dependencies, configures middleware, registers routes, and starts the server.
func main() {
	var (
		injector = do.New()
	)

	provider.RegisterDependencies(injector)

	if !args(injector) {
		return
	}

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.RegisterRoutes(server, injector)

	run(server)
}
