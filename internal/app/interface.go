package app

import (
	"github.com/ss7labs/amts26/config"
	"gorm.io/gorm"
)

type AppContext interface {
	DB() *gorm.DB
	Config() *config.AppConfig
}
