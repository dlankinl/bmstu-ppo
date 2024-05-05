package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewLogger() (logger *zap.SugaredLogger, err error) {
	err = os.Mkdir("./logs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("создание папки для логов: %w", err)
	}

	var files = []struct {
		Path  string
		Level zapcore.LevelEnabler
	}{
		{
			Path: "./logs/error",
			Level: zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level == zapcore.ErrorLevel
			}),
		},
		{
			Path: "./logs",
			Level: zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level == zapcore.InfoLevel
			}),
		},
	}

	var cores []zapcore.Core
	for _, f := range files {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   f.Path,
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     40,
				Compress:   true,
				LocalTime:  true,
			}),
			f.Level,
		))
	}

	cores = append(cores, zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zapcore.InfoLevel,
	))

	loggr := zap.New(zapcore.NewTee(cores...))

	return loggr.Sugar(), nil
}
