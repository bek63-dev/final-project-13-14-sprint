// Package main запускает HTTP-сервер планировщика задач:
// собирает зависимости (конфигурацию, подключение к БД, маршрутизатор) и
// передаёт управление серверу.
package main

import (
	"log"
	"os"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/config"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/router"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/server"

	_ "modernc.org/sqlite"
)

func main() {
	logger := log.New(os.Stdout, "SERVER: ", log.Ldate|log.Ltime)

	// Загрузка конфигурации из .env и переменных окружения.
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("config.Load: %w", err)
	}

	// Подключение к БД и применение схемы.
	db, err := database.New(cfg.DBFile)
	if err != nil {
		logger.Fatalf("database.New: %v", err)
	}
	defer db.Close()

	// Сборка маршрутизатора со всеми HTTP-обработчиками.
	mux := router.New(db, cfg)

	// Запуск сервера, блокирует выполнение до остановки.
	srv := server.New(cfg.Address(), logger, mux)
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
