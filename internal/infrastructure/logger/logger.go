package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(level string) *Logger {
	var zapLogger *zap.Logger
	var err error

	if level == "production" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	sugar := zapLogger.Sugar()
	return &Logger{sugar}
}

func (l *Logger) Sync() {
	_ = l.SugaredLogger.Sync()
}
