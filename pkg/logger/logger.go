package logger

import (
    "os"
    "time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// Config holds the logger configuration
type Config struct {
    Level      string
    TimeFormat string
    Pretty     bool
    Caller     bool
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() Config {
    return Config{
        Level:      "info",
        TimeFormat: time.RFC3339,
        Pretty:     true,
        Caller:     true,
    }
}

// Init initializes the logger with the given configuration
func Init(config Config) {
    // Set up output
    var output zerolog.ConsoleWriter
    if config.Pretty {
        output = zerolog.ConsoleWriter{
            Out:        os.Stdout,
            TimeFormat: config.TimeFormat,
            NoColor:    false,
        }
    } else {
        output = zerolog.ConsoleWriter{
            Out:        os.Stdout,
            TimeFormat: config.TimeFormat,
            NoColor:    true,
        }
    }

    // Create logger
    logger := zerolog.New(output).
        With().
        Timestamp()

    // Add caller information if enabled
    if config.Caller {
        logger = logger.Caller()
    }

    // Set global logger
    log.Logger = logger.Logger()

    // Set global level
    SetLogLevel(config.Level)

    // Log initialization
    log.Info().
        Str("level", config.Level).
        Bool("pretty", config.Pretty).
        Bool("caller", config.Caller).
        Msg("Logger initialized")
}

// GetLogger returns the global logger instance
func GetLogger() *zerolog.Logger {
    return &log.Logger
}

// SetLogLevel sets the global log level
func SetLogLevel(level string) {
    switch level {
    case "debug":
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
        log.Debug().Msg("Log level set to debug")
    case "info":
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
        log.Info().Msg("Log level set to info")
    case "warn":
        zerolog.SetGlobalLevel(zerolog.WarnLevel)
        log.Warn().Msg("Log level set to warn")
    case "error":
        zerolog.SetGlobalLevel(zerolog.ErrorLevel)
        log.Error().Msg("Log level set to error")
    default:
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
        log.Warn().Str("level", level).Msg("Invalid log level, defaulting to info")
    }
}

// Log levels
const (
    LevelDebug = "debug"
    LevelInfo  = "info"
    LevelWarn  = "warn"
    LevelError = "error"
)

// Common log messages
const (
    MsgServerStarted     = "Server started"
    MsgServerStopped     = "Server stopped"
    MsgDatabaseConnected = "Database connected"
    MsgConfigLoaded      = "Configuration loaded"
    MsgRequestReceived   = "Request received"
    MsgRequestProcessed  = "Request processed"
    MsgErrorOccurred     = "Error occurred"
)

// Log fields
const (
    FieldRequestID    = "request_id"
    FieldClientIP     = "client_ip"
    FieldMethod       = "method"
    FieldPath         = "path"
    FieldStatusCode   = "status_code"
    FieldLatency      = "latency"
    FieldErrorMessage = "error_message"
    FieldErrorType    = "error_type"
    FieldErrorDetails = "error_details"
) 