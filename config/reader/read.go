package reader

import (
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	viperDefaultDelimiter = "."
)

func Read(config interface{}, opts ...viper.DecoderConfigOption) error {
	viperLogger := viper.New()
	viperLogger.SetEnvKeyReplacer(strings.NewReplacer(viperDefaultDelimiter, "_")) // replace default viper delimiter for env vars
	viperLogger.AutomaticEnv()
	viperLogger.SetTypeByDefaultValue(true)

	defaults.SetDefaults(viperLogger)
	err := viperLogger.Unmarshal(config, opts...)
	if err != nil {
		return errors.WithMessage(err, "failed to parse configuration")
	}

	return nil
}
