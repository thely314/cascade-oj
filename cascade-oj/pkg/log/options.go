package log

// Options contains the configuration options for the logger.
type Options struct {
	Level         string // Log level (e.g., "INFO", "WARN", "ERROR")
	Filename      string // Log filename
	MaxSize       int    // Maximum size of a single log file in megabytes
	MaxBackups    int    // Maximum number of old log files to retain
	MaxAge        int    // Maximum number of days to retain log files
	EnableConsole bool   // Whether to output to stdout
}

// Option is a function that modifies the Options.
type Option func(*Options)

// WithLevel sets the log level.
//
// valid levels are DEBUG, INFO, WARN, ERROR, FATAL
func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithFilename sets the log filename.
func WithFilename(filename string) Option {
	return func(o *Options) {
		o.Filename = filename
	}
}

// WithMaxSize sets the maximum size of a single log file.
func WithMaxSize(size int) Option {
	return func(o *Options) {
		o.MaxSize = size
	}
}

// WithMaxBackups sets the maximum number of old log files to retain.
func WithMaxBackups(backups int) Option {
	return func(o *Options) {
		o.MaxBackups = backups
	}
}

// WithMaxAge sets the maximum number of days to retain log files.
func WithMaxAge(age int) Option {
	return func(o *Options) {
		o.MaxAge = age
	}
}

// WithEnableConsole enables or disables console output.
func WithEnableConsole(enable bool) Option {
	return func(o *Options) {
		o.EnableConsole = enable
	}
}
