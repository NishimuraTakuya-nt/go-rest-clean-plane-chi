package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var Config appConfig

type appConfig struct {
	Env            string        `mapstructure:"env" validate:"required"`
	LogLevel       string        `mapstructure:"log_level"`
	Version        string        `mapstructure:"version"` // fixme increment version
	ServerAddress  string        `mapstructure:"server_address" validate:"required"`
	AllowedOrigins []string      `mapstructure:"allowed_origins" validate:"required"`
	JWTSecretKey   string        `mapstructure:"jwt_secret_key" validate:"required"`
	RequestTimeout time.Duration `mapstructure:"request_timeout" validate:"required"`
	// OpenTelemetry
	ServiceName   string        `mapstructure:"service_name" validate:"required"`
	DDAgentHost   string        `mapstructure:"dd_agent_host" validate:"required"`
	DDAgentPort   string        `mapstructure:"dd_agent_port" validate:"required"`
	SamplingRate  float64       `mapstructure:"sampling_rate" validate:"required,min=0,max=1"`
	BatchTimeout  time.Duration `mapstructure:"batch_timeout" validate:"required"`
	BatchSize     int           `mapstructure:"batch_size" validate:"required"`
	SamplingRatio float64       `mapstructure:"sampling_ratio" validate:"required"`
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

	// OpenTelemetry
	v.SetDefault("service_name", "go-rest-clean-plane-chi")
	v.SetDefault("dd_agent_host", "localhost")
	v.SetDefault("dd_agent_port", "8126")
	v.SetDefault("sampling_rate", 1.0)
	v.SetDefault("batch_timeout", 1*time.Second)
	v.SetDefault("batch_size", 512)
	v.SetDefault("sampling_ratio", 1.0)

	viper.AutomaticEnv()
	if err := v.Unmarshal(&Config); err != nil {
		panic(err)
	}
}

// Validate validates the config values.
func (c *appConfig) Validate() error {
	return validator.New().Struct(c)
}
