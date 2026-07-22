// Package database отвечает за подключение к SQLite и работу с таблицей scheduler
package database

import (
	"database/sql"
	"fmt"
	"log"
)

// schema описывает структуру таблицы scheduler и индекс по колонке date
const schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(128) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);`

// DB оборачивает стандартный *sql.DB, чтобы можно было навешивать на него собственные методы
type DB struct {
	*sql.DB
}

// New открывает соединение с БД по указанному пути, проверяет его через Ping
// и применяет schema (создание таблицы/индекса, если их ещё нет)
func New(dbFile string) (*DB, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	log.Printf("Connecting to database %q", dbFile)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	log.Printf("Initializing database %q", dbFile)
	if _, err = db.Exec(schema); err != nil {
		return nil, fmt.Errorf("db.Exec: %w", err)
	}

	return &DB{db}, nil
}
