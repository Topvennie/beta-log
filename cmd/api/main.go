package main

import (
	"context"
	"fmt"

	"github.com/Topvennie/beta-log/internal/climb"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server"
	"github.com/Topvennie/beta-log/internal/task"
	"github.com/Topvennie/beta-log/pkg/config"
	"github.com/Topvennie/beta-log/pkg/db"
	"github.com/Topvennie/beta-log/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	loggerFile := config.GetString("logger.file")
	loggerLevelStr := config.GetString("logger.level")
	var loggerLevel *zapcore.Level
	if loggerLevelStr != "" {
		loggerLevelTmp, err := zapcore.ParseLevel(loggerLevelStr)
		if err != nil {
			panic(fmt.Errorf("invalid logger level %s | %w", loggerLevelStr, err))
		}
		loggerLevel = &loggerLevelTmp
	}

	zapLogger, err := logger.New(logger.Config{
		Console: true,
		File:    loggerFile,
		Level:   loggerLevel,
	})
	if err != nil {
		panic(fmt.Errorf("zap logger initialization failed: %w", err))
	}
	zap.ReplaceGlobals(zapLogger)

	db, err := db.NewPSQL()
	if err != nil {
		zap.S().Fatalf("Unable to connect to database %v", err)
	}

	repository.Init(db)

	if err := task.Init(); err != nil {
		zap.S().Fatalf("Failed to init task %v", err)
	}
	if err := climb.New().Start(context.Background()); err != nil {
		zap.S().Fatalf("Failed to start climb %v", err)
	}

	s, err := server.New()
	if err != nil {
		zap.S().Fatalf("Failed to create the server %v", err)
	}

	zap.S().Infof("Server is running on %s", s.Addr)
	if err := s.Listen(s.Addr); err != nil {
		zap.S().Fatalf("Failure while running the server %v", err)
	}
}
