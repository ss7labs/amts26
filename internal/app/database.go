package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ss7labs/amts26/config"
	"github.com/ss7labs/amts26/pkg/common"

	sloggorm "github.com/imdatngo/slog-gorm/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func getDatabase(cfg config.DBConfig, logger *slog.Logger) *gorm.DB {
	return getPgDatabase(cfg, logger)
}

func getPgDatabase(config config.DBConfig, slogger *slog.Logger) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Ashgabat",
		config.Host,
		config.User,
		config.Passwd,
		config.Name,
		config.Port)
	cfg := sloggorm.NewConfig(slogger.Handler()).
		WithGroupKey("db").
		WithSlowThreshold(200 * time.Millisecond) // Flag queries taking longer than 200ms

	glogger := sloggorm.NewWithConfig(cfg)
	pool, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},

		Logger: glogger,
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
