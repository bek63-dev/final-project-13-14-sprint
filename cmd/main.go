// Точка входа приложения.
// Собирает зависимости (конфиг, БД, роутер, сервер) и запускает HTTP-сервер планировщика задач.
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
	logger := log.New(os.Stdout, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	cfg := config.Load()

	db, err := database.New(cfg.DBFile)
	if err != nil {
		logger.Fatalf("database.New: %v", err)
	}
	defer db.Close()

	mux := router.New(db)

	srv := server.New(cfg.Address(), logger, mux)
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}