package log

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitLog init zap log
func InitLog() (*zap.Logger, error) {
	var cfg zap.Config

	if disableJSON := viper.GetBool("log.disable_json"); disableJSON {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.OutputPaths = []string{
		viper.GetString("log.file"),
	}

	return cfg.Build()
}
