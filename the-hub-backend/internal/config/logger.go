package config

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	var logger *zap.Logger
	var err error

	if env == "production" {
		// Log rotation using lumberjack
		lumberjackLogger := &lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    50, // MB
			MaxBackups: 7,
			MaxAge:     30,   // days
			Compress:   true, // gzip old files
		}

		writeSyncer := zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),        // optional console logging
			zapcore.AddSync(lumberjackLogger), // file logging
		)

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			writeSyncer,
			zap.InfoLevel,
		)

		logger = zap.New(core, zap.AddCaller())
	} else {
		// Dev: coloured console logs
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = cfg.Build(zap.AddCaller())
		if err != nil {
			panic("Failed to initialize logger: " + err.Error())
		}
	}

	Logger = logger.Sugar()
}

// EnsureDir creates a directory and all necessary parents
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// SaveBase64Image saves a base64 encoded image to a file
func SaveBase64Image(base64Data, filePath string) error {
	// Remove data URL prefix if present
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	// Decode base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := EnsureDir(filepath.Dir(filePath)); err != nil {
		return err
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the image data
	_, err = file.Write(imageData)
	return err
}

// DeleteFile deletes a file if it exists
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}
	return os.Remove(filePath)
}
