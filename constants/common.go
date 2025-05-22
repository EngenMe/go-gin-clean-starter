package constants

const (

	// ENUM_ROLE_ADMIN represents the constant value for the admin role within the system.
	ENUM_ROLE_ADMIN = "admin"

	// ENUM_ROLE_USER represents the constant value assigned to the "user" role within the system.
	ENUM_ROLE_USER = "user"

	// ENUM_RUN_DEVELOPMENT represents the environment mode for development configurations and settings.
	ENUM_RUN_DEVELOPMENT = "dev"

	// ENUM_RUN_PRODUCTION represents the environment mode for production configurations and settings.
	ENUM_RUN_PRODUCTION = "prod"

	// ENUM_RUN_TESTING represents the environment mode for testing configurations and settings.
	ENUM_RUN_TESTING = "test"

	// ENUM_PAGINATION_PER_PAGE defines the default number of items to display per page in pagination.
	ENUM_PAGINATION_PER_PAGE = 10

	// ENUM_PAGINATION_PAGE represents the default starting page number used for pagination.
	ENUM_PAGINATION_PAGE = 1

	// DB is a constant key used to identify and provide the main database dependency in the dependency injection container.
	DB = "db"

	// JWTService is a constant key used for identifying the JWT service dependency in the dependency injection container.
	JWTService = "JWTService"
)
