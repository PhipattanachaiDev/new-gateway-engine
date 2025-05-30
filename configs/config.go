package configs

import (
	"fmt"
	models "naturelink/models/config_models"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return err
	}
	return nil
}

func GetFromEnv() (string, string, error) {
	ip := os.Getenv("PUBLIC_IP")
	port := os.Getenv("PORT")

	if ip == "" || port == "" {
		return "", "", fmt.Errorf("missing PUBLIC_IP or PORT in environment variables")
	}

	return ip, port, nil
}

func GetSizeLog() (models.Logger, error) {
	maxsize := os.Getenv("LOG_MAX_SIZE")
	maxbackups := os.Getenv("LOG_MAX_BACKUPS")
	maxage := os.Getenv("LOG_MAX_AGE")
	compressENV := os.Getenv("LOG_COMPRESS")
	logLevel := os.Getenv("LOG_LEVEL")

	compress, err := strconv.ParseBool(compressENV)
	if err != nil {
		compress = true
	}

	return models.Logger{
		MaxSize:    maxsize,
		MaxBackups: maxbackups,
		MaxAge:     maxage,
		Compress:   compress,
		Level:      logLevel,
	}, nil
}

func GetServiceLogPath() (models.ServiceLog, error) {
	logPath := os.Getenv("LOG_PATH")
	logFileName := os.Getenv("LOG_FILE")
	serviceName := os.Getenv("SERVICE_NAME")

	return models.ServiceLog{
		LogPath:     logPath,
		LogFile:     logFileName,
		ServiceName: serviceName,
	}, nil
}
