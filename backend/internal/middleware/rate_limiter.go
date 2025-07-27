package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimiterConfig holds configuration for the rate limiter
type RateLimiterConfig struct {
	Limit         int           // Number of requests allowed in the window
	Window        time.Duration // Time window for rate limiting
	BlockDuration time.Duration // How long to block after exceeding limit
}

// DefaultRateLimiterConfig provides sensible defaults
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Limit:         100,             // 100 requests
		Window:        time.Minute,     // per minute
		BlockDuration: time.Minute * 5, // block for 5 minutes
	}
}

// NewRateLimiter creates a new rate limiter with default configuration
func NewRateLimiter(logger *zap.SugaredLogger) gin.HandlerFunc {
	return RateLimiter(DefaultRateLimiterConfig(), logger)
}

// RateLimiter implements a rate limiting middleware with IP-based tracking
func RateLimiter(config RateLimiterConfig, logger *zap.SugaredLogger) gin.HandlerFunc {
	type client struct {
		count        int
		lastReset    time.Time
		blockedUntil time.Time
	}

	var (
		clients = make(map[string]*client)
		mux     sync.Mutex
	)

	// Start a cleanup goroutine to prevent memory leaks
	go func() {
		for {
			time.Sleep(time.Hour) // Cleanup once per hour
			now := time.Now()
			mux.Lock()
			for ip, cli := range clients {
				// Remove entries that haven't been seen in a while and aren't blocked
				if now.Sub(cli.lastReset) > time.Hour && now.After(cli.blockedUntil) {
					delete(clients, ip)
				}
			}
			mux.Unlock()
		}
	}()

	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		mux.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{count: 0, lastReset: time.Now()}
		}

		cli := clients[ip]
		now := time.Now()

		// Check if client is currently blocked
		if now.Before(cli.blockedUntil) {
			remaining := cli.blockedUntil.Sub(now).Seconds()
			mux.Unlock()
			logger.Warnw("Blocked client attempted request",
				"ip", ip,
				"blocked_for", remaining,
			)
			c.Header("Retry-After", fmt.Sprintf("%.0f", remaining))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: "Too many requests, please try again later",
				},
			})
			return
		}

		// Reset counter if window has passed
		if time.Since(cli.lastReset) > config.Window {
			cli.count = 0
			cli.lastReset = now
		}

		// Increment request counter
		cli.count++

		// Check if over limit
		if cli.count > config.Limit {
			// Block the client
			cli.blockedUntil = now.Add(config.BlockDuration)
			mux.Unlock()
			logger.Warnw("Rate limit exceeded, client blocked",
				"ip", ip,
				"count", cli.count,
				"limit", config.Limit,
				"block_duration", config.BlockDuration,
			)
			c.Header("Retry-After", fmt.Sprintf("%.0f", config.BlockDuration.Seconds()))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{
				Success: false,
				Error: &ErrorData{
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: "Rate limit exceeded, please try again later",
				},
			})
			return
		}

		mux.Unlock()
		c.Next()
	}
}
