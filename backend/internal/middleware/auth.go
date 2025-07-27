package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JWTConfig struct {
	Secret            string
	TokenExpiration   time.Duration
	RefreshExpiration time.Duration
	Issuer            string
	Audience          []string
}

// DefaultJWTConfig returns a default JWT configuration
func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:            "your-secret-key",  // This should be overridden with a real secret
		TokenExpiration:   time.Hour * 24,     // 24 hours
		RefreshExpiration: time.Hour * 24 * 7, // 7 days
		Issuer:            "lightweight-crm",
		Audience:          []string{"api"},
	}
}

// JWTClaims represents custom JWT claims
type JWTClaims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWT(config JWTConfig, logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    CodeUnauthorized,
					Message: "Authorization header is missing",
				},
			})
			return
		}

		// validate the authorization format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    CodeUnauthorized,
					Message: "Invalid authorization header format",
				},
			})
			return
		}

		tokenString := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// validate the token signature
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Secret), nil
		})

		if err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    CodeUnauthorized,
					Message: "Invalid token",
				},
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    CodeUnauthorized,
					Message: "Invalid token",
				},
			})
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Next()
	}
}
