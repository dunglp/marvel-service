package service

import (
	"io"
	"strings"

	"github.com/rs/zerolog"
)

var (
	logLevels = map[string]zerolog.Level{
		"TRACE": zerolog.TraceLevel,
		"DEBUG": zerolog.DebugLevel,
		"INFO":  zerolog.InfoLevel,
		"WARN":  zerolog.WarnLevel,
		"ERROR": zerolog.ErrorLevel,
		"OFF":   zerolog.NoLevel,
	}
)

func parseLogLevel(level string) zerolog.Level {
	logLevel := zerolog.NoLevel
	if val, ok := logLevels[strings.ToUpper(level)]; ok {
		logLevel = val
	}
	return logLevel
}

func setLogLevelFieldName(fieldName string) {
	zerolog.LevelFieldName = fieldName
}

func setLogLevel(level string) {
	zerolog.SetGlobalLevel(parseLogLevel(level))
}

func newLogger(output io.Writer) zerolog.Logger {
	return zerolog.New(output).With().Timestamp().Logger()
}
