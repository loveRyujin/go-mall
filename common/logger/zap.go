package logger

import (
	"github.com/loveRyujin/go-mall/common/enum"
	"github.com/loveRyujin/go-mall/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func init() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	fileWriteSyncer := getFileLogWriter()

	var cores []zapcore.Core
	switch config.App.Env {
	case enum.ModeTest, enum.ModeProd:
		cores = append(cores, zapcore.NewCore(encoder, fileWriteSyncer, zapcore.InfoLevel))
	case enum.ModeDev:
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
		cores = append(cores, zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel))
	}
	core := zapcore.NewTee(cores...)
	_logger = zap.New(core)
}

var _logger *zap.Logger

func getFileLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  config.App.Log.FilePath,
		MaxSize:   config.App.Log.FileMaxSize,
		MaxAge:    config.App.Log.FileMaxAge,
		Compress:  false,
		LocalTime: true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
