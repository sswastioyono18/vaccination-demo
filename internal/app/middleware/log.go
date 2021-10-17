package middleware

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var Logger *zap.Logger

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func NewLogger(env string)  {
	var cfg zap.Config
	fmt.Println("Using",env,"config")
	if env == "DEV" {
		cfg = zap.NewDevelopmentConfig()
	} else if env == "PROD" {
		cfg = zap.NewProductionConfig()
	}

	cfg.EncoderConfig.EncodeTime = syslogTimeEncoder
	_, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}