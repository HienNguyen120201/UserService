package util

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(lv zapcore.Level, pretty bool) (*zap.Logger, error) {
	c := zap.NewDevelopmentConfig()
	var opts []zap.Option
	if pretty {
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	}
	level := zap.NewAtomicLevel()

	if err := level.UnmarshalText([]byte(lv.String())); err != nil {
		return nil, fmt.Errorf("could not parse log level %s", lv.String())
	}
	c.Level = level
	return c.Build(opts...)
}
