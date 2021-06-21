package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	EncodingTypeConsole = "console"
	EncodingTypeJson    = "json"
)

func CreateLogger(logLevel Level, encoding string, fields ...zap.Field) (*zap.Logger, error) {
	encoderConfig := defaultEncoderConfiguration(encoding)

	factory := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Encoding:         encoding,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig:    encoderConfig,
	}

	return factory.Build(zap.Fields(fields...))
}

// DefaultConfiguration applies the default settings of the logger configuration
func defaultEncoderConfiguration(encodingType string) zapcore.EncoderConfig {
	var encoderConfig = zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "loglevel",
		TimeKey:      "@timestamp",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "calling_func",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	switch encodingType {
	case EncodingTypeConsole:
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case EncodingTypeJson:
		fallthrough
	default:
		encodingType = EncodingTypeJson
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	return encoderConfig
}
