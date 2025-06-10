package middleware

import (
    "bytes"
    "io"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
    w.body.Write(b)
    return w.ResponseWriter.Write(b)
}

// Logger is a middleware that logs request/response
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()

        // Read request body
        var requestBody []byte
        if c.Request.Body != nil {
            requestBody, _ = io.ReadAll(c.Request.Body)
            c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
        }

        // Create custom response writer
        blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
        c.Writer = blw

        // Process request
        c.Next()

        // Calculate latency
        latency := time.Since(start)

        // Get response body
        responseBody := blw.body.String()

        // Log request details
        log.Info().
            Str("request_id", c.GetString("request_id")).
            Str("client_ip", c.ClientIP()).
            Str("method", c.Request.Method).
            Str("path", c.Request.URL.Path).
            Str("query", c.Request.URL.RawQuery).
            Str("user_agent", c.Request.UserAgent()).
            Str("request_body", string(requestBody)).
            Int("status_code", c.Writer.Status()).
            Str("response_body", responseBody).
            Dur("latency", latency).
            Int("body_size", c.Writer.Size()).
            Str("error", c.Errors.ByType(gin.ErrorTypePrivate).String()).
            Msg("Request processed")

        // Log additional metrics
        log.Info().
            Str("request_id", c.GetString("request_id")).
            Dur("latency", latency).
            Int("status_code", c.Writer.Status()).
            Msg("Request metrics")
    }
}

// RequestID is a middleware that adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get request ID from header or generate new one
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }

        // Set request ID in context
        c.Set("request_id", requestID)

        // Add request ID to response header
        c.Header("X-Request-ID", requestID)

        c.Next()
    }
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
    return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
} 