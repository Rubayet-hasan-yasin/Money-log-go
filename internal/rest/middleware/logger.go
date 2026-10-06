package middleware

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.buf.Len() < 500 {
		remaining := 500 - w.buf.Len()
		if len(b) > remaining {
			w.buf.Write(b[:remaining])
		} else {
			w.buf.Write(b)
		}
	}
	return w.ResponseWriter.Write(b)
}

func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		shouldCaptureRequest := method == "POST" || method == "PUT" || method == "PATCH"
		if shouldCaptureRequest && c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		rw := &responseWriter{ResponseWriter: c.Writer, buf: bufferPool.Get().(*bytes.Buffer)}
		c.Writer = rw

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		shouldLog := statusCode >= 400 || shouldCaptureRequest || method == "GET"
		if !shouldLog {
			rw.buf.Reset()
			bufferPool.Put(rw.buf)
			return
		}

		logFields := make([]zap.Field, 0, 8)
		logFields = append(logFields,
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
		)

		if query := c.Request.URL.RawQuery; query != "" {
			logFields = append(logFields, zap.String("query", query))
		}

		if clientIP := c.ClientIP(); clientIP != "127.0.0.1" && clientIP != "::1" {
			logFields = append(logFields, zap.String("ip", clientIP))
		}

		log.Info("http request", logFields...)

		rw.buf.Reset()
		bufferPool.Put(rw.buf)
	}
}
