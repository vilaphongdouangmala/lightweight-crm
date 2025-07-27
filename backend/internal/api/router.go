package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vilaphongdouangmala/lightweight-crm/backend/internal/config"
	"github.com/vilaphongdouangmala/lightweight-crm/backend/internal/middleware"
	"go.uber.org/zap"
)

func SetupRouter(cfg *config.Config, logger *zap.SugaredLogger) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.NewRateLimiter(logger))
	router.Use(middleware.CORS())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/")
		SetupPublicRoutes(public, cfg, logger)

		// Create JWT config
		jwtConfig := middleware.DefaultJWTConfig()
		jwtConfig.Secret = cfg.Auth.JWTSecret

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.JWT(jwtConfig, logger))
		SetupProtectedRoutes(protected, cfg, logger)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// SetupPublicRoutes configures the public routes
func SetupPublicRoutes(router *gin.RouterGroup, cfg *config.Config, logger *zap.SugaredLogger) {
	// // Add auth controller routes (login, register)
	// authController := NewAuthController(cfg, logger)
	// router.POST("/auth/login", authController.Login)
	// router.POST("/auth/register", authController.Register)
}

func SetupProtectedRoutes(router *gin.RouterGroup, cfg *config.Config, logger *zap.SugaredLogger) {
	// // Add user controller routes
	//
	//	userController := NewUserController(cfg, logger)
	//	router.GET("/users/me", userController.GetCurrentUser)
	//	router.PUT("/users/me", userController.UpdateCurrentUser)
}
