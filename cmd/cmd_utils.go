package cmd

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// getLogger returns a zap.Logger instance and an error.
//
// It retrieves the log level from the environment variable LOG_LEVEL. If the
// variable is empty, it defaults to "info". It then creates a zap.Config
// instance based on whether the PRODUCTION_MODE environment variable is set to
// "true". If it is, the config is set to production mode with output paths set
// to "stdout" and "stderr", and encoding set to either "console" or "json". If
// PRODUCTION_MODE is not set to "true", the config is set to development mode
// and the encoder is configured to use colorized output. The log level is then
// parsed and if there is an error, an error is returned. Finally, the function
// builds the logger and returns it.
//
// Return Values:
// - *zap.Logger: The zap.Logger instance.
// - error: An error if there was an error parsing the log level.
func getLogger(logLevel string, productionMode bool) (*zap.Logger, error) {
	if logLevel == "" {
		logLevel = "info"
	}

	var cfg zap.Config

	if productionMode {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		cfg.Encoding = "console" // "console" | "json"

	} else {
		cfg = zap.NewDevelopmentConfig()
		// Use only with console encoder (i.e. not in production)
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, fmt.Errorf("getLogger unmarshal level %q: %v", logLevel, err)
	}
	cfg.Level = level
	if err != nil {
		return nil, fmt.Errorf("getLogger: %v", err)
	}

	return cfg.Build()
}
