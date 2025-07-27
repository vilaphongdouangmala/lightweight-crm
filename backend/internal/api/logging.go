package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes a production-ready logger
func InitLogger() (*zap.Logger, error) {
	// Define log level based on environment
	logLevel := zapcore.InfoLevel
	if gin.Mode() == gin.DebugMode {
		logLevel = zapcore.DebugLevel
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		logLevel,
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// RequestLogger logs detailed information about each request
func RequestLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request details
		logger.Infow("Request received",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		// Process request
		c.Next()

		// Log response details
		logger.Infow("Response sent",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
			"errors", c.Errors.Errors(),
		)
	}
}
