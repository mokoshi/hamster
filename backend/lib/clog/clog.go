package clog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO 環境によって振り分ける
var config = zap.Config{
	Level:    zap.NewAtomicLevelAt(zapcore.InfoLevel),
	Encoding: "console",
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:        "Time",
		LevelKey:       "Level",
		NameKey:        "Name",
		CallerKey:      "Caller",
		MessageKey:     "Msg",
		StacktraceKey:  "St",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	},
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

//var logger, _ = zap.NewProduction()
var logger, _ = config.Build()

var Logger = logger.Sugar()
