package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel represents logging severity.
type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

// String returns the string representation of a log level.
func (l LogLevel) String() string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides simple logging functionality.
type Logger struct {
	level  LogLevel
	output io.Writer
}

// DefaultLogger is the global logger instance.
var DefaultLogger = NewLogger(LogInfo)

// NewLogger creates a new logger with the specified level.
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		output: os.Stdout,
	}
}

// SetLevel changes the logging level.
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput changes the logging output destination.
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// log writes a log message if the level is high enough.
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "[%s] [%s] %s\n", timestamp, level.String(), message)
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogDebug, format, args...)
}

// Info logs an info message.
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogInfo, format, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogWarn, format, args...)
}

// Error logs an error message.
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogError, format, args...)
}

// Package-level convenience functions

// Debug logs a debug message using the default logger.
func Debug(format string, args ...interface{}) {
	DefaultLogger.Debug(format, args...)
}

// Info logs an info message using the default logger.
func Info(format string, args ...interface{}) {
	DefaultLogger.Info(format, args...)
}

// Warn logs a warning message using the default logger.
func Warn(format string, args ...interface{}) {
	DefaultLogger.Warn(format, args...)
}

// Error logs an error message using the default logger.
func Error(format string, args ...interface{}) {
	DefaultLogger.Error(format, args...)
}

// SetLogLevel sets the level for the default logger.
func SetLogLevel(level LogLevel) {
	DefaultLogger.SetLevel(level)
}

// ParseLogLevel converts a string to LogLevel.
func ParseLogLevel(s string) LogLevel {
	switch s {
	case "debug", "DEBUG":
		return LogDebug
	case "info", "INFO":
		return LogInfo
	case "warn", "WARN", "warning", "WARNING":
		return LogWarn
	case "error", "ERROR":
		return LogError
	default:
		return LogInfo
	}
}
