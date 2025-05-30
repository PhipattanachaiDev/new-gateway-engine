package logservices

import (
	// "ezviewlite/configs"
	"fmt"
	"naturelink/configs"
	"os"
	"strconv"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var (
	logger  *zerolog.Logger
	logFile *lumberjack.Logger
)

// InitLogger initializes the logger with the specified log file name.
func InitLogger(pathName string, serviceName string, logLevel zerolog.Level) {
	// Load logger config from environment
	loggerConfig, err := configs.GetSizeLog()
	if err != nil {
		fmt.Println("Error loading logger config:", err)
		return
	}

	// Convert string values to int
	maxSize, err := strconv.Atoi(loggerConfig.MaxSize)
	if err != nil {
		fmt.Println("Invalid MaxSize value:", err)
		return
	}

	maxBackups, err := strconv.Atoi(loggerConfig.MaxBackups)
	if err != nil {
		fmt.Println("Invalid MaxBackups value:", err)
		return
	}

	maxAge, err := strconv.Atoi(loggerConfig.MaxAge)
	if err != nil {
		fmt.Println("Invalid MaxAge value:", err)
		return
	}

	// Configure log rotation
	logFile = &lumberjack.Logger{
		Filename:   pathName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   loggerConfig.Compress,
	}

	// Set the global log level
	zerolog.SetGlobalLevel(logLevel)

	// Initialize the global logger
	log := zerolog.New(logFile).Level(logLevel).With().Str("serviceName", serviceName).Timestamp().Logger()
	logger = &log // Assign logger as a pointer
}

// GetLogger returns the global logger instance.
func GetLogger() *zerolog.Logger {
	if logger == nil {
		defaultLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger = &defaultLogger
	}
	return logger
}

// CloseLogger should be called when the application is shutting down
// to close the log file properly.
func CloseLogger() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			fmt.Printf("Error closing log file: %s\n", err)
		}
	}
}