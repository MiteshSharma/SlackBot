package logger

import "github.com/MiteshSharma/SlackBot/config"

type TestLogger struct {
}

func NewTestLogger(loggerParam config.LoggerConfig) *TestLogger {
	testLogger := &TestLogger{}
	return testLogger
}

func (zl *TestLogger) Debug(message string, args ...Argument) {
}

func (zl *TestLogger) Info(message string, args ...Argument) {
}

func (zl *TestLogger) Warn(message string, args ...Argument) {
}

func (zl *TestLogger) Error(message string, args ...Argument) {
}
