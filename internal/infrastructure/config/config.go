package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var Config AppConfig

type AppConfig struct {
	Env            string        `mapstructure:"env" validate:"required"`
	LogLevel       string        `mapstructure:"log_level"`
	Version        string        `mapstructure:"version"` // fixme increment version
	ServerAddress  string        `mapstructure:"server_address" validate:"required"`
	AllowedOrigins []string      `mapstructure:"allowed_origins" validate:"required"`
	JWTSecretKey   string        `mapstructure:"jwt_secret_key" validate:"required"`
	RequestTimeout time.Duration `mapstructure:"request_timeout" validate:"required"`
	// DataDog Agent
	DDAgentHost string `mapstructure:"dd_agent_host" validate:"required"`
	DDAgentPort string `mapstructure:"dd_agent_port" validate:"required"`
}

func init() {
	v := viper.New()
	v.SetDefault("env", "dev")
	v.SetDefault("log_level", "INFO")
	v.SetDefault("server_address", ":8081")
	v.SetDefault("allowed_origins", []string{"*"})
	v.SetDefault("jwt_secret_key", "jwt-secret")
	v.SetDefault("request_timeout", 180*time.Second)
	//v.SetDefault("request_timeout", 1*time.Second) // fixme

	// DataDog Agent
	v.SetDefault("dd_agent_host", "localhost")
	v.SetDefault("dd_agent_port", "4317")

	v.AutomaticEnv()
	if err := v.Unmarshal(&Config); err != nil {
		panic(err)
	}
}

// Validate validates the config values.
func (c *AppConfig) Validate() error {
	return validator.New().Struct(c)
}
