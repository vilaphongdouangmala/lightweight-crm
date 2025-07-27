package middleware

import (
	"fmt"
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

func GenerateToken(userId, role string, config JWTConfig) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Audience:  config.Audience,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(config.TokenExpiration)),
			Subject:   userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(userId string, config JWTConfig) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    config.Issuer,
		Subject:   userId,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(config.RefreshExpiration)),
		ID:        fmt.Sprintf("refresh_%d", now.Unix()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    CodeForbidden,
					Message: "You do not have permission to access this resource",
				},
			})
			return
		}

		roleString, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    CodeForbidden,
					Message: "You do not have permission to access this resource",
				},
			})
			return
		}

		hasRole := false
		for _, role := range roles {
			if roleString == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    CodeForbidden,
					Message: "You do not have permission to access this resource",
				},
			})
			return
		}

		c.Next()
	}
}
