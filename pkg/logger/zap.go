// Package logger initiates a zap logger
package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
		if _, err := os.Stat(logCfg.File); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("checking if log file exists %+v | %w", logCfg, err)
			}

			dir := filepath.Dir(logCfg.File)
			if dir != "" {
				// nolint:gosec // Let me set my own permissions
				if err := os.MkdirAll(dir, 0o0755); err != nil {
					return nil, fmt.Errorf("creating log directory %+v | %w", logCfg, err)
				}
			}

			// nolint:gosec // Let me set my own permissions
			if err := os.WriteFile(logCfg.File, nil, 0o0644); err != nil {
				return nil, fmt.Errorf("creating log file %+v | %w", logCfg, err)
			}
		}
	}

	outputPaths := []string{}
	errorOutputPaths := []string{}
	if logCfg.File != "" {
		outputPaths = append(outputPaths, logCfg.File)
		errorOutputPaths = append(errorOutputPaths, logCfg.File)
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
