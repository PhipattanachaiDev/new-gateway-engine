package logservices

import (
	"fmt"
	"os"
)

func isDebugMode() bool {
	var debugModeEnv = "DEBUG_MODE"
	return os.Getenv(debugModeEnv) == "true"
}

func ServicesInfo(message interface{}) {
	var msg string

	// แปลง message เป็น string
	switch v := message.(type) {
	case string:
		msg = v
	case int, float64:
		msg = fmt.Sprintf("%v", v)
	default:
		msg = "unsupported type"
	}

	logMessage("info", msg)
}

func ServicesError(err error) {
	logMessage("error", err.Error())
}

func ServicesWarning(err error) {
	logMessage("warn", err.Error())
}

func ServicesFatal(err error) {
	logMessage("fatal", err.Error())
}

func ServicesDebug(message string) {
	logMessage("debug", message)
}

func logMessage(level, message string) {
	if !isDebugMode() {
		return
	}

	logger := GetLogger()
	switch level {
	case "info":
		logger.Info().Msg(message)
	case "error":
		logger.Error().Msg(message)
	case "warn":
		logger.Warn().Msg(message)
	case "fatal":
		logger.Fatal().Msg(message)
	case "debug":
		logger.Debug().Msg(message)
	}
}

