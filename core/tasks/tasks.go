package tasks

import (
	"backend/core/enviroment"
	"backend/core/internal/database"
	"sync"

	"github.com/rs/zerolog/log"
)

type appData struct {
	database *database.Database
	settings *enviroment.Environment
}

var (
	instance *appData
	once     sync.Once
)

// OnStartup инициализирует все зависимости при запуске.
func OnStartup() {
	log.Debug().Msg("Running on startup tasks")

	once.Do(func() {
		settings := enviroment.GetInstance()
		db := database.GetInstance(settings.DSN())

		instance = &appData{
			database: db,
			settings: settings,
		}
	})
}

// OnShutdown закрывает подключения при завершении.
func OnShutdown() {
	log.Debug().Msg("Running shutdown tasks")
	if instance != nil && instance.database != nil {
		instance.database.Close()
	}
}
