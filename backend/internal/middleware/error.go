package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool       `json:"success"`
	Error   *ErrorData `json:"error"`
}

// ErrorData contains detailed error information
type ErrorData struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// CustomError represents an application error with status code
type CustomError struct {
	Code       string
	Message    string
	Details    interface{}
	StatusCode int
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return e.Message
}

// ErrorHandler handles errors
func ErrorHandler(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Check if it's our custom error type
				if customErr, ok := e.Err.(*CustomError); ok {
					c.JSON(customErr.StatusCode, ErrorResponse{
						Success: false,
						Error: &ErrorData{
							Code:    customErr.Code,
							Message: customErr.Message,
							Details: customErr.Details,
						},
					})
					logger.Errorw("Application error",
						"error", customErr.Message,
						"code", customErr.Code,
						"status", customErr.StatusCode,
					)
				} else {
					// Handle generic errors
					c.JSON(http.StatusInternalServerError, ErrorResponse{
						Success: false,
						Error: &ErrorData{
							Code:    "INTERNAL_SERVER_ERROR",
							Message: "An internal server error occurred",
						},
					})
					logger.Errorw("Unexpected error",
						"error", e.Err,
						"path", c.Request.URL.Path,
						"method", c.Request.Method,
					)
				}
				for _, e := range c.Errors {
					logger.Error(e)
					break
				}
			}
		}
	}
}

// Error constants
const (
	CodeBadRequest          = "BAD_REQUEST"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeForbidden           = "FORBIDDEN"
	CodeNotFound            = "NOT_FOUND"
	CodeConflict            = "CONFLICT"
	CodeValidationFailed    = "VALIDATION_FAILED"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string, details interface{}) *CustomError {
	return &CustomError{
		Code:       CodeBadRequest,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *CustomError {
	return &CustomError{
		Code:       CodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *CustomError {
	return &CustomError{
		Code:       CodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *CustomError {
	return &CustomError{
		Code:       CodeNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// NewInternalServerError creates an internal server error
func NewInternalServerError(message string) *CustomError {
	return &CustomError{
		Code:       CodeInternalServerError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}
