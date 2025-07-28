package repository

import (
	"fmt"
	"time"

	"github.com/vilaphongdouangmala/lightweight-crm/backend/internal/config"
	"github.com/vilaphongdouangmala/lightweight-crm/backend/internal/domain"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB     *gorm.DB
	logger *zap.SugaredLogger
}

func NewDatabase(cfg *config.Config, zapLogger *zap.SugaredLogger) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)
	// Configure GORM logger
	gormLogLevel := logger.Silent
	if cfg.Server.Mode == "debug" {
		gormLogLevel = logger.Info
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC() // Use UTC for all timestamps
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	return &Database{
		DB:     db,
		logger: zapLogger,
	}, nil
}

// AutoMigrate runs database migrations
func (d *Database) AutoMigrate() error {
	d.logger.Info("Running database migrations")

	err := d.DB.AutoMigrate(
		&domain.User{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	d.logger.Info("Database migrations completed successfully")
	return nil
}
