package app

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ss7labs/amts26/config"
	"github.com/ss7labs/amts26/pkg/common"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func getDatabase(cfg config.DBConfig) *gorm.DB {
	return getPgDatabase(cfg)
}

func getPgDatabase(config config.DBConfig) *gorm.DB {
	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host,
		config.User,
		config.Passwd,
		config.Name,
		config.Port)
	pool, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.NewSlogLogger(slogLogger, // io writer
			logger.Config{
				SlowThreshold:             time.Millisecond * 200,                                                // Slow SQL threshold
				LogLevel:                  common.If(config.Debug, logger.Info, logger.Silent).(logger.LogLevel), //nolint:errcheck // type assertion is safe
				IgnoreRecordNotFoundError: true,                                                                  // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,                                                                 // Disable color
			}),
	})

	common.Must(err)
	sqlDB, err := pool.DB()
	common.Must(err)
	// SetMaxIdleConns sets the maximum number of idle connections in the pool
	sqlDB.SetMaxIdleConns(config.IdleConn)
	// SetMaxOpenConns sets the maximum number of open database connections
	sqlDB.SetMaxOpenConns(config.MaxConn)
	// SetConnMaxLifetime sets the maximum lifetime a connection can be reused
	// sqlDB.SetConnMaxLifetime(time.Hour * 8)
	return pool
}
