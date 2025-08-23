package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

// Level represents the logging level
type Level int

const (
	// DebugLevel is the most verbose logging level
	DebugLevel Level = iota
	// InfoLevel is for informational messages
	InfoLevel
	// WarnLevel is for warning messages
	WarnLevel
	// ErrorLevel is for error messages
	ErrorLevel
	// FatalLevel is for fatal errors that cause the program to exit
	FatalLevel
)

// Logger provides structured logging with levels
type Logger struct {
	level    Level
	output   io.Writer
	useColor bool
	prefix   string

	// Color functions
	debugColor *color.Color
	infoColor  *color.Color
	warnColor  *color.Color
	errorColor *color.Color
	fatalColor *color.Color
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(InfoLevel, os.Stderr, true)
}

// New creates a new logger with the specified level and output
func New(level Level, output io.Writer, useColor bool) *Logger {
	return &Logger{
		level:      level,
		output:     output,
		useColor:   useColor,
		debugColor: color.New(color.FgCyan),
		infoColor:  color.New(color.FgGreen),
		warnColor:  color.New(color.FgYellow),
		errorColor: color.New(color.FgRed),
		fatalColor: color.New(color.FgRed, color.Bold),
	}
}

// SetLevel sets the global logging level
func SetLevel(level Level) {
	defaultLogger.level = level
}

// SetVerbose enables or disables verbose logging
func SetVerbose(verbose bool) {
	if verbose {
		defaultLogger.level = DebugLevel
	} else {
		defaultLogger.level = InfoLevel
	}
}

// SetOutput sets the output writer for the default logger
func SetOutput(w io.Writer) {
	defaultLogger.output = w
}

// SetColorEnabled enables or disables color output
func SetColorEnabled(enabled bool) {
	defaultLogger.useColor = enabled
	// Also force fatih/color global to honor this toggle
	color.NoColor = !enabled
}

// formatMessage formats a log message with timestamp and level
func (l *Logger) formatMessage(level string, msg string) string {
	timestamp := time.Now().Format("15:04:05")
	if l.prefix != "" {
		return fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, l.prefix, level, msg)
	}
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, msg)
}

// log writes a log message at the specified level
func (l *Logger) log(level Level, levelStr string, colorFunc *color.Color, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	msg := fmt.Sprintf(format, args...)
	formattedMsg := l.formatMessage(levelStr, msg)

	if l.useColor && colorFunc != nil {
		colorFunc.Fprintln(l.output, formattedMsg)
	} else {
		fmt.Fprintln(l.output, formattedMsg)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, "DEBUG", l.debugColor, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, "INFO", l.infoColor, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, "WARN", l.warnColor, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, "ERROR", l.errorColor, format, args...)
}

// Fatal logs a fatal error message and exits the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, "FATAL", l.fatalColor, format, args...)
	os.Exit(1)
}

// Global logging functions using the default logger

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// Fatal logs a fatal error message and exits the program
func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// LogError logs a structured error with appropriate detail level
func LogError(err error, verbose bool) {
	if err == nil {
		return
	}

	if verbose {
		// In verbose mode, show full error details
		Error("Error occurred: %+v", err)
	} else {
		// In normal mode, show just the message
		Error("%v", err)
	}
}

// ParseLevel parses a string level into a Level type
func ParseLevel(levelStr string) Level {
	switch levelStr {
	case "debug", "DEBUG":
		return DebugLevel
	case "info", "INFO":
		return InfoLevel
	case "warn", "WARN", "warning", "WARNING":
		return WarnLevel
	case "error", "ERROR":
		return ErrorLevel
	case "fatal", "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}
