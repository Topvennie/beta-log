// Package logger initiates a zap logger
package logger

import (
	"errors"
	"fmt"
	"os"

	"github.com/Topvennie/beta-log/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	File    string
	Console bool
	Level   *zapcore.Level
}

func New(logCfg Config) (*zap.Logger, error) {
	if logCfg.File != "" {
		// nolint:gosec // I want to set my own permissions
		err := os.Mkdir("logs", 0o755)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("create logs directory %w", err)
		}
	}

	outputPaths := []string{}
	errorOutputPaths := []string{}
	if logCfg.File != "" {
		outputPaths = append(outputPaths, fmt.Sprintf("logs/%s.log", logCfg.File))
		errorOutputPaths = append(errorOutputPaths, fmt.Sprintf("logs/%s.log", logCfg.File))
	}
	if logCfg.Console {
		outputPaths = append(outputPaths, "stdout")
		errorOutputPaths = append(errorOutputPaths, "stderr")
	}

	if len(outputPaths) == 0 {
		return nil, errors.New("no output paths specified")
	}

	var cfg zap.Config

	if config.IsDev() {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")

		level := zap.DebugLevel
		if logCfg.Level != nil {
			level = *logCfg.Level
		}
		cfg.Level.SetLevel(level)
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")

		level := zap.WarnLevel
		if logCfg.Level != nil {
			level = *logCfg.Level
		}
		cfg.Level.SetLevel(level)
	}

	cfg.OutputPaths = outputPaths
	cfg.ErrorOutputPaths = errorOutputPaths

	logger, err := cfg.Build(zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		return nil, err
	}

	return logger, nil
}
