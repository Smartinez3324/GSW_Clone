package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// init Initializes the logger
// The logger is configured using the logger.yaml file in the data/config directory
// If the file is not found, the logger will default to a development logger
func init() {
	//Default Logger
	defaultLogger := zap.Must(zap.NewDevelopment())

	// Viper config parsing
	viperConfig := viper.New()
	viperConfig.SetConfigType("yaml")
	viperConfig.SetConfigName("logger")
	viperConfig.AddConfigPath("data/config")

	if err := viperConfig.ReadInConfig(); err != nil {
		defaultLogger.Warn(fmt.Sprint(err))
		logger = defaultLogger
		return
	}
	// Create zap config
	var loggerConfig zap.Config

	outputPaths := viperConfig.GetStringSlice("OutputPaths")

	// Make and populate the output paths
	for index, path := range outputPaths {
		if path == "stdout" || path == "stderr" {
			outputPaths[index] = path
			continue
		}
		// Create log file
		logFileName := fmt.Sprint("gsw_service_log-", time.Now().Format("2006-01-02 15:04:05"), ".log")
		baseName := fmt.Sprint(path, logFileName)
		totalLogPath := baseName

		// Ensures unique file name
		numIncrease := 1

		for {
			if _, err := os.Stat(totalLogPath); os.IsNotExist(err) {
				break
			}
			totalLogPath = fmt.Sprintf("%s.%d", baseName, numIncrease)
			numIncrease++
		}

		_, noPath := os.Create(totalLogPath)
		if noPath != nil {
			os.Mkdir(path, 0755)
			os.Create(totalLogPath)
		}
		outputPaths[index] = totalLogPath
	}
	// Setting Logger Paths
	loggerConfig.OutputPaths = outputPaths
	loggerConfig.ErrorOutputPaths = outputPaths

	// Setting Logger Level
	level, err := zap.ParseAtomicLevel(viperConfig.GetString("level"))

	if err != nil {
		defaultLogger.Warn(fmt.Sprint(err))
	}
	loggerConfig.Level = level

	// Setting Encoding Type
	var levelEncoder zapcore.LevelEncoder
	var timeEncoder zapcore.TimeEncoder
	var durationEncoder zapcore.DurationEncoder
	var callerEncoder zapcore.CallerEncoder

	levelEncoder.UnmarshalText([]byte(viperConfig.GetString("encoderConfig.levelEncoder")))
	timeEncoder.UnmarshalText([]byte(viperConfig.GetString("encoderConfig.timeEncoder")))
	durationEncoder.UnmarshalText([]byte(viperConfig.GetString("encoderConfig.durationEncoder")))
	callerEncoder.UnmarshalText([]byte(viperConfig.GetString("encoderConfig.callerEncoder")))

	loggerConfig.Encoding = viperConfig.GetString("encoding")
	loggerConfig.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     viperConfig.GetString("encoderConfig.messageKey"),
		LevelKey:       viperConfig.GetString("encoderConfig.levelKey"),
		TimeKey:        viperConfig.GetString("encoderConfig.timeKey"),
		NameKey:        viperConfig.GetString("encoderConfig.nameKey"),
		CallerKey:      viperConfig.GetString("encoderConfig.callerKey"),
		StacktraceKey:  viperConfig.GetString("encoderConfig.stacktraceKey"),
		LineEnding:     viperConfig.GetString("encoderConfig.LineEnding"),
		EncodeLevel:    levelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: durationEncoder,
		EncodeCaller:   callerEncoder,
	}

	logger = zap.Must(loggerConfig.Build())
}

// Info logs an info message
func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

// Warn logs a warning message
func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

// Debug logs a debug message
func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

// Fatal logs a fatal message
func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

// Error logs an error message
func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

// Panic logs a panic message
func Panic(message string, fields ...zap.Field) {
	logger.Panic(message, fields...)
}
