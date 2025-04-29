// config/config.go (adding to existing config)
package config

type (
	Config struct {
		//Server   ServerConfig
		//Database DatabaseConfig
		JWT JWTConfig
		// Other existing config
	}

	JWTConfig struct {
		SecretKey          string `mapstructure:"SECRET_KEY"`
		Issuer             string `mapstructure:"ISSUER"`
		RefreshTokenExpiry int    `mapstructure:"REFRESH_TOKEN_EXPIRY"` // in hours
		RotateRefreshToken bool   `mapstructure:"ROTATE_REFRESH_TOKEN"`
	}
)

// Init function would need to be updated to load these new config values
