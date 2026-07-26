// Package database отвечает за подключение к SQLite и хранение задач
// планировщика в таблице scheduler. Используется обработчиками пакета
// api как единственная точка доступа к данным.
package database

import (
	"database/sql"
	"fmt"
	"log"
)

// schema описывает структуру таблицы scheduler и индекс по колонке date,
// применяется при инициализации БД.
const schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(128) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);`

// DB оборачивает стандартный *sql.DB, позволяя навешивать на него
// собственные методы работы с задачами (см. task.go).
type DB struct {
	*sql.DB
}

// New открывает соединение с БД по пути dbFile, проверяет его через Ping
// и применяет schema, создавая таблицу и индекс, если их ещё нет.
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
