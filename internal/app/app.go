package app

import (
	"log/slog"
	"os"

	"github.com/ss7labs/amts26/config"
	"gorm.io/gorm"
)

type Application struct {
	appConfig *config.AppConfig
	gormDB    *gorm.DB
}

var (
	_ AppContext = (*Application)(nil)
)

func NewApplication(cfg *config.AppConfig) *Application {
	app := &Application{
		appConfig: cfg,
	}
	return app
}
func (a *Application) DB() *gorm.DB {
	return a.gormDB
}
func (a *Application) Config() *config.AppConfig {
	return a.appConfig
}

func (a *Application) Init() {
	slog.SetDefault(a.setLogger())
	a.gormDB = getDatabase(a.appConfig.Database)
	slog.Info("DB driver loaded", "DB type", a.appConfig.Database.Type)
}
func (a *Application) setLogger() *slog.Logger {
	var logLevel slog.Level
	if a.appConfig.System.Debug {
		logLevel = slog.LevelDebug.Level()
	}
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	// 2. Create the handler (slog.NewTextHandler or slog.NewJSONHandler)
	handler := slog.NewTextHandler(os.Stdout, opts)
	// 3. Set it as the application-wide default logger
	logger := slog.New(handler)
	return logger
}
