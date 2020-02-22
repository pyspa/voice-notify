package log

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger        *zap.Logger
	commonOptions = []zap.Option{
		zap.AddCallerSkip(1),
	}
)

type Fields map[string]interface{}

func Init() error {
	outPath := viper.GetString("log.log")
	errPath := viper.GetString("log.err_log")

	l, err := buildLogger(outPath, errPath)
	if err != nil {
		return errors.Wrap(err, "failed build logger")
	}
	logger = l

	return nil
}

func buildLogger(outPath, errPath string) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.InfoLevel)

	debug := viper.GetBool("log.debug")
	if debug {
		level.SetLevel(zap.DebugLevel)
	}

	cfg := zap.Config{
		Level:       level,
		Development: true,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		DisableCaller:    false,
		OutputPaths:      []string{outPath},
		ErrorOutputPaths: []string{errPath},
	}

	return cfg.Build()
}

func Debug(msg string, f Fields) {
	logger.WithOptions(commonOptions...).WithOptions(commonOptions...).Debug(msg, zapValues(f)...)
}

func Info(msg string, f Fields) {
	logger.WithOptions(commonOptions...).Info(msg, zapValues(f)...)
}

func Error(msg interface{}, f Fields) {
	err, ok := msg.(error)
	if ok {
		logger.WithOptions(commonOptions...).Error(fmt.Sprintf("%+v", err), zapValues(f)...)
	} else {
		logger.WithOptions(commonOptions...).Error(fmt.Sprintf("%v", msg), zapValues(f)...)
	}
}

func zapValues(f Fields) (fs []zap.Field) {
	for k, v := range f {
		fs = append(fs, zap.Any(k, v))
	}
	return
}

func Close() {
	_ = logger.Sync()
}
