package provider

import (
	"github.com/samber/do"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/controller"
	"github.com/Caknoooo/go-gin-clean-starter/repository"
	"github.com/Caknoooo/go-gin-clean-starter/service"
)

// ProvideUserDependencies initializes and provides user-related dependencies, including repositories, services, and controllers.
var ProvideUserDependencies = func(injector *do.Injector) {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)

	userService := service.NewUserService(userRepository, refreshTokenRepository, jwtService, db)

	do.Provide(
		injector, func(i *do.Injector) (controller.UserController, error) {
			return controller.NewUserController(userService), nil
		},
	)
}
