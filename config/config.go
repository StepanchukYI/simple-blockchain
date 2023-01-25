package config

import (
	"github.com/StepanchukYI/simple-blockchain/server/http"
)

type Config struct {
	LogLevel   string      `mapstructure:"LOG_LEVEL" default:"DEBUG"`
	HTTP       http.Config `mapstructure:"HTTP_PORT"  default:"8080"`
	WSSPort    int         `mapstructure:"WSS_PORT"  default:"8081"`
	HealthPort int         `mapstructure:"HEALTH_PORT"  default:"8082"`
}
