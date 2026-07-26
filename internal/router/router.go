// Package router собирает http.ServeMux и регистрирует в нём все HTTP-маршруты приложения,
// связывая их с обработчиками пакета api и раздачей статики фронтенда.
package router

import (
	"net/http"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/api"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/config"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
)

// webDir — путь к директории со статическими файлами фронтенда (index.html, js, css)
const webDir = "./web"

// New создаёт и настраивает http.ServeMux со всеми маршрутами приложения.
// Вызывается в main при сборке зависимостей сервера.
func New(db *database.DB, cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	Init(mux, db, cfg)
	return mux
}

// Init регистрирует маршруты API и статики в mux. Создаёт api.Handler,
// связывая его с БД и конфигурацией, и оборачивает маршруты задач
// middleware api.Auth для проверки JWT-токена.
func Init(mux *http.ServeMux, db *database.DB, cfg *config.Config) {
	handler := api.NewHandler(db, cfg)

	// Раздача статических файлов фронтенда.
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	// Публичные маршруты, не требующие аутентификации.
	mux.HandleFunc("/api/nextdate", api.NextDayHandler)
	mux.HandleFunc("/api/signin", handler.SignInHandler)

	// Маршруты задач, защищённые проверкой JWT-токена.
	mux.HandleFunc("/api/task", handler.Auth(handler.TaskHandler))
	mux.HandleFunc("/api/tasks", handler.Auth(handler.GetTasksHandler))
	mux.HandleFunc("/api/task/done", handler.Auth(handler.DoneTaskHandler))
}
