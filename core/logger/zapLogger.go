package logger

import (
	"github.com/MiteshSharma/SlackBot/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	Zap *zap.Logger
}

func NewLogger(loggerParam config.LoggerConfig) *ZapLogger {
	zapConfig := generateConfig(loggerParam)
	logger, _ := zapConfig.Build(zap.AddCallerSkip(1), zap.AddCaller())
	zapLogger := &ZapLogger{
		Zap: logger,
	}
	return zapLogger
}

func generateConfig(loggerParam config.LoggerConfig) zap.Config {
	loggerConfig := zap.NewProductionConfig()
	if (config.LoggerConfig{}) != loggerParam {
		loggerConfig.Encoding = "json"
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		if loggerParam.LogFilePath != "" {
			loggerConfig.OutputPaths = []string{"stderr", loggerParam.LogFilePath}
			loggerConfig.ErrorOutputPaths = []string{"stderr", loggerParam.LogFilePath}
		}
		loggerConfig.EncoderConfig = zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		}
	}
	return loggerConfig
}

func (zl *ZapLogger) Debug(message string, args ...Argument) {
	zl.Zap.Debug(message, args...)
}

func (zl *ZapLogger) Info(message string, args ...Argument) {
	zl.Zap.Info(message, args...)
}

func (zl *ZapLogger) Warn(message string, args ...Argument) {
	zl.Zap.Warn(message, args...)
}

func (zl *ZapLogger) Error(message string, args ...Argument) {
	zl.Zap.Error(message, args...)
}
