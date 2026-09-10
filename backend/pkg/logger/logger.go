package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(mode string) (*zap.Logger, error) {
	var config zapcore.EncoderConfig
	var encoder zapcore.Encoder
	var level zapcore.Level

	if mode == "release" {
		config = zap.NewProductionEncoderConfig()
		config.TimeKey = "timestamp"
		config.EncodeTime = zapcore.ISO8601TimeEncoder

		encoder = zapcore.NewJSONEncoder(config)

		level = zapcore.InfoLevel
	} else {
		config = zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")

		encoder = zapcore.NewConsoleEncoder(config)

		level = zapcore.DebugLevel
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level,)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}