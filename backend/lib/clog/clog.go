package clog

import "go.uber.org/zap"

// TODO 環境によって振り分ける
//var logger, _ = zap.NewProduction()
var logger, _ = zap.NewDevelopment()

var Logger = logger.Sugar()
