package database

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// Database — singleton-обёртка над подключением к PostgreSQL.
type Database struct {
	Conn *pgx.Conn
}

var (
	dbInstance *Database
	dbOnce     sync.Once
)

// GetInstance возвращает единственное подключение к БД (создаёт при первом вызове).
func GetInstance(dsn string) *Database {
	dbOnce.Do(func() {
		conn, err := pgx.Connect(context.Background(), dsn)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to connect to database")
		}

		if err := conn.Ping(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Unable to ping database")
		}

		log.Info().Msg("Successfully connected to database")
		dbInstance = &Database{Conn: conn}
	})
	return dbInstance
}

// Close закрывает подключение к БД.
func (d *Database) Close() {
	if d.Conn != nil {
		d.Conn.Close(context.Background())
	}
}
