package app

import (
	"github.com/ss7labs/amts26/config"
	"github.com/ss7labs/amts26/pkg/rtprefix"
	"gorm.io/gorm"
)

type AppContext interface {
	DB() *gorm.DB
	Config() *config.AppConfig
	GetPrefixStore() *rtprefix.PrefixStore
}
