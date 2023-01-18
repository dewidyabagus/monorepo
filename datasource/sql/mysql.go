package sql

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLConfig struct {
	Host             string
	Port             int
	User             string
	Password         string
	Schema           string
	MaxIdleConns     int
	MaxOpenConns     int
	ConnMaxLifetime  int
	Environment      string
	SlowLogThreshold int
}

func NewMySQL(config *MySQLConfig) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Schema,
	)

	if config.SlowLogThreshold == 0 {
		config.SlowLogThreshold = 200 // Default: 200 Milisecond
	}

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Duration(config.SlowLogThreshold) * time.Millisecond,
				IgnoreRecordNotFoundError: true,
			},
		),
	}

	if strings.TrimSpace(strings.ToLower(config.Environment)) == "prod" {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				// Slow SQL threshold
				SlowThreshold:             time.Duration(config.SlowLogThreshold) * time.Millisecond,
				LogLevel:                  logger.Error, // Log level
				IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,        // Disable color
			},
		)
	}

	if db, err = gorm.Open(mysql.Open(dsn), gormConfig); err != nil {
		return nil, fmt.Errorf("open database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("instance db config failed: %w", err)
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.ConnMaxLifetime))

	return
}
