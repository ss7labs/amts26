package config

import (
	"os"

	"github.com/ss7labs/amts26/pkg/common"
	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	Type     string `yaml:"type"`      // Database type: postgres or sqlite
	Host     string `yaml:"host"`      // PostgreSQL host address
	Port     int    `yaml:"port"`      // PostgreSQL port
	Name     string `yaml:"name"`      // Database name or SQLite file path
	User     string `yaml:"user"`      // PostgreSQL username
	Passwd   string `yaml:"passwd"`    // PostgreSQL password
	MaxConn  int    `yaml:"max_conn"`  // Maximum connections
	IdleConn int    `yaml:"idle_conn"` // Idle connections
	Debug    bool   `yaml:"debug"`     // Debug mode
}

type SysConfig struct {
	Debug bool `yaml:"debug"`
}

type LogConfig struct {
	Mode       string `yaml:"mode"`
	FileEnable bool   `yaml:"file_enable"`
	Filename   string `yaml:"filename"`
}

type AppConfig struct {
	System   SysConfig `yaml:"system" json:"system"`
	Database DBConfig  `yaml:"database" json:"database"`
	Logger   LogConfig `yaml:"logger" json:"logger"`
}

func LoadConfig(cfile string) *AppConfig {
	if cfile == "" {
		cfile = "config.yml"
	}

	cfg := new(AppConfig)
	if common.FileExists(cfile) {
		data := common.Must2(os.ReadFile(cfile))        //nolint:gosec // G304: config file path is intentionally variable
		common.Must(yaml.Unmarshal(data.([]byte), cfg)) //nolint:errcheck // type assertion is safe after Must2
	} else {
		cfg = DefaultAppConfig
	}

	return cfg
}

var DefaultAppConfig = &AppConfig{
	System: SysConfig{
		Debug: true,
	},
	Database: DBConfig{
		Type:     "postgres",  // Default to SQLite for development and testing
		Host:     "127.0.0.1", // PostgreSQL configuration (used when type is postgres)
		Port:     5433,
		Name:     "amts", // SQLite: database filename; PostgreSQL: database name
		User:     "postgres",
		Passwd:   "pwd",
		MaxConn:  100,
		IdleConn: 10,
		Debug:    false,
	},
	Logger: LogConfig{
		Mode:       "development",
		FileEnable: true,
		Filename:   "/tmp/amts.log",
	},
}
