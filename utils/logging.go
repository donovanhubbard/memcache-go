package utils

import (
	"encoding/json"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func getLogLevel() zapcore.Level {
	defaultLevel := zapcore.FatalLevel

	envLevel := os.Getenv("LOG_LEVEL")

	if envLevel == "" {
		return defaultLevel
	}
	if strings.ToLower(envLevel) == "debug" {
		return zapcore.DebugLevel
	}
	if strings.ToLower(envLevel) == "info" {
		return zapcore.InfoLevel
	}
	if strings.ToLower(envLevel) == "warn" {
		return zapcore.WarnLevel
	}
	if strings.ToLower(envLevel) == "error" {
		return zapcore.ErrorLevel
	}
	if strings.ToLower(envLevel) == "panic" {
		return zapcore.PanicLevel
	}
	return defaultLevel
}

func InitializeLogger() {

	rawJSON := []byte(`{
		"level": "fatal",
		"encoding": "json",
		"outputPaths": ["stderr"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	  }`)
  
	  var cfg zap.Config
	  if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		  panic(err)
	  }

	  cfg.Level.SetLevel(getLogLevel())

	  cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	  cfg.EncoderConfig.TimeKey = "timestamp"

	  Logger := zap.Must(cfg.Build())
	  defer Logger.Sync()

	  Sugar = Logger.Sugar()
  
	  Sugar.Info("Logging initialized")
}