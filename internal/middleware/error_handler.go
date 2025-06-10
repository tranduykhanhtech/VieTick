package middleware

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "vietick/pkg/errors"
)

// ErrorResponse represents the error response structure
type ErrorResponse struct {
    Timestamp string      `json:"timestamp"`
    RequestID string      `json:"request_id"`
    Type      string      `json:"type"`
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Details   string      `json:"details,omitempty"`
    Path      string      `json:"path"`
    Method    string      `json:"method"`
}

// ErrorHandler is a middleware that handles errors
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // Check if there are any errors
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            var response ErrorResponse

            // Set common fields
            response.Timestamp = time.Now().Format(time.RFC3339)
            response.RequestID = c.GetString("request_id")
            response.Path = c.Request.URL.Path
            response.Method = c.Request.Method

            // Log the error
            log.Error().
                Str("request_id", response.RequestID).
                Str("path", response.Path).
                Str("method", response.Method).
                Err(err).
                Msg("Request error")

            switch e := err.(type) {
            case *errors.AppError:
                response.Type = string(e.Type)
                response.Code = e.Code
                response.Message = e.Message
                response.Details = e.Details

                // Log specific error details
                log.Error().
                    Str("type", response.Type).
                    Int("code", response.Code).
                    Str("message", response.Message).
                    Str("details", response.Details).
                    Msg("Application error")
            default:
                response.Type = string(errors.ErrorTypeInternal)
                response.Code = http.StatusInternalServerError
                response.Message = "Internal server error"
                response.Details = "An unexpected error occurred"

                // Log internal error
                log.Error().
                    Err(err).
                    Msg("Internal server error")
            }

            c.JSON(response.Code, response)
        }
    }
}

// RecoveryHandler is a middleware that recovers from panics
func RecoveryHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log the panic
                log.Error().
                    Interface("error", err).
                    Str("path", c.Request.URL.Path).
                    Str("method", c.Request.Method).
                    Msg("Panic recovered")

                // Create error response
                response := ErrorResponse{
                    Timestamp: time.Now().Format(time.RFC3339),
                    RequestID: c.GetString("request_id"),
                    Type:      string(errors.ErrorTypeInternal),
                    Code:      http.StatusInternalServerError,
                    Message:   "Internal server error",
                    Details:   "The server encountered an unexpected condition",
                    Path:      c.Request.URL.Path,
                    Method:    c.Request.Method,
                }

                c.JSON(response.Code, response)
                c.Abort()
            }
        }()
        c.Next()
    }
} 