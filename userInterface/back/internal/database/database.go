package database

import (
	"log"
	"os"
	"time"

	"github.com/rosty-git/Smartchain-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config interface {
	GetEnv() string
	GetDsn() string
	GetGormDebug() bool
}

func New(dsn string, debug bool) (*gorm.DB, func() error, error) {
	var gormConfig *gorm.Config

	if debug {
		gormConfig = &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // You can customize the logger output
				logger.Config{
					SlowThreshold: time.Second, // SQL queries that take longer than this threshold will be logged as slow queries
					LogLevel:      logger.Info, // Set log level to Log mode to log all queries
					Colorful:      true,        // Enable colorful output
				},
			),
		}
	} else {
		gormConfig = &gorm.Config{}
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return db, sqlDb.Close, nil
}

func Initialize(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Account{})

	return err
}
