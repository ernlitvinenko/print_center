package storage

import (
	"backend/core/enviroment"
	"backend/core/internal/database"
	"backend/core/repositories"

	"github.com/jackc/pgx/v5"
)

// Storage — базовый интерфейс для всех хранилищ.
type Storage interface{}

type storage struct {
	conn     *pgx.Conn
	queries  *repositories.Queries
	settings *enviroment.Environment
}

// GetInstance создаёт экземпляр storage с доступом к БД и настройкам.
func GetInstance[T Storage]() T {
	db := database.GetInstance(enviroment.GetInstance().DSN())
	queries := repositories.New(db.Conn)
	settings := enviroment.GetInstance()

	instance := &storage{
		conn:     db.Conn,
		queries:  queries,
		settings: settings,
	}

	result, ok := any(instance).(T)
	if !ok {
		panic("storage instance does not implement requested interface")
	}
	return result
}
