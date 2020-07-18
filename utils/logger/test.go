package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type LevelEnabler bool

func (o LevelEnabler) Enabled(level zapcore.Level) bool {
	return bool(o)
}

func GetTestLogger() *zap.Logger {
	core, _ := observer.New(LevelEnabler(false))
	return zap.New(core)
}

func InitTest() {
	Log = Logger{Logger: GetTestLogger()}
}