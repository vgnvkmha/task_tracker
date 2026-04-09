package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Encoding = "console"
	config.EncoderConfig.StacktraceKey = ""

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()

	return sugar, nil
}
