package log

import (
	"io"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger creates a new logger instance.
func NewLogger(opts ...Option) (log.Logger, error) {
	// default Options
	options := Options{
		Level:         "INFO", // support DEBUG, INFO, WARN, ERROR, FATAL
		EnableConsole: true,
		MaxSize:       100, // 100MB
		MaxBackups:    10,
		MaxAge:        7, // 7 days
	}

	// Apply all incoming Option functions
	for _, o := range opts {
		o(&options)
	}

	var writers []io.Writer

	// Set up file rotation (with lumberjack) if a filename is provided
	if options.Filename != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   options.Filename,
			MaxSize:    options.MaxSize,
			MaxBackups: options.MaxBackups,
			MaxAge:     options.MaxAge,
			LocalTime:  true,
		}
		writers = append(writers, fileWriter)
	}

	// Add console writer if enabled
	if options.EnableConsole {
		writers = append(writers, os.Stdout)
	}
	multiWriter := io.MultiWriter(writers...)

	// Create Kratos logger
	logger := log.NewStdLogger(multiWriter)
	level := log.ParseLevel(options.Level)

	return log.NewFilter(logger, log.FilterLevel(level)), nil
}
