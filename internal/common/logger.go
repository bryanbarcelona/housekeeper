// internal/common/logger.go
package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Global loggers
var (
	Info  = &timestampedLogger{out: os.Stdout, prefix: "[INFO] "}
	Warn  = &timestampedLogger{out: os.Stdout, prefix: "[WARN] "}
	Error = &timestampedLogger{out: os.Stderr, prefix: "[ERROR] "}
	Debug = &timestampedLogger{out: os.Stdout, prefix: "[DEBUG] "}

	isDebugEnabled bool
	logFile        *os.File
)

// SetDebug enables or disables debug logs
func SetDebug(enabled bool) {
	isDebugEnabled = enabled
}

func CloseLogFile() error {
	if logFile != nil {
		f := logFile
		logFile = nil
		return f.Close()
	}
	return nil
}

func initLogger(cfg LoggingConfig) error {
	var writer io.Writer = os.Stdout

	if cfg.LogToFile {
		logDir := filepath.Dir(cfg.LogFilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		f, err := os.OpenFile(cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		logFile = f

		if cfg.AlsoPrintToConsole {
			writer = io.MultiWriter(f, os.Stdout)
		} else {
			writer = f
		}
	} else {
		if logFile != nil {
			logFile.Close()
			logFile = nil
		}
	}

	Info.out = writer
	Warn.out = writer
	Error.out = writer
	Debug.out = writer

	SetDebug(cfg.Debug)
	return nil
}

/*
func initLogger(cfg LoggingConfig) error {
	var writer io.Writer = os.Stdout

	if cfg.LogToFile {
		logDir := filepath.Dir(cfg.LogFilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		f, err := os.OpenFile(cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		if cfg.AlsoPrintToConsole {
			writer = io.MultiWriter(f, os.Stdout)
		} else {
			writer = f
		}
	}

	// Reassign output writers
	Info.out = writer
	Warn.out = writer
	Error.out = writer
	Debug.out = writer

	SetDebug(cfg.Debug)
	return nil
}
*/

// SetupLogging initializes logging system
func SetupLogging(cfg LoggingConfig) error {
	if cfg.LogToFile {
		if err := RotateLog(cfg.LogFilePath, 5); err != nil {
			return err
		}
	}

	return initLogger(cfg)
}

// timestampedLogger wraps log output with timestamps
type timestampedLogger struct {
	out    io.Writer
	prefix string
}

func (l *timestampedLogger) Println(v ...interface{}) {
	if strings.HasPrefix(l.prefix, "[DEBUG]") && !isDebugEnabled {
		return // skip if debug disabled
	}

	timestamp := time.Now().Format(time.RFC3339)
	msg := fmt.Sprint(v...)
	fmt.Fprintf(l.out, "[%s] %s%s\n", timestamp, l.prefix, msg)
}

func (l *timestampedLogger) Printf(format string, v ...interface{}) {
	if strings.HasPrefix(l.prefix, "[DEBUG]") && !isDebugEnabled {
		return // skip if debug disabled
	}

	timestamp := time.Now().Format(time.RFC3339)
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(l.out, "[%s] %s%s\n", timestamp, l.prefix, msg)
}
