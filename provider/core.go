package provider

import (
	"github.com/samber/do"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/config"
	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/service"
)

// InitDatabase configures and registers the main database instance into the dependency injection container.
func InitDatabase(injector *do.Injector) {
	do.ProvideNamed(
		injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
			return config.SetUpDatabaseConnection(), nil
		},
	)
}

// RegisterDependencies encapsulates the registration of all necessary application dependencies into the dependency injector.
var RegisterDependencies = func(injector *do.Injector) {
	InitDatabase(injector)

	do.ProvideNamed(
		injector, constants.JWTService, func(i *do.Injector) (service.JWTService, error) {
			return service.NewJWTService(), nil
		},
	)

	ProvideUserDependencies(injector)
}
