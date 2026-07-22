// Package router отвечает за создание http.ServeMux и регистрацию в нём всех обработчиков
package router

import (
	"net/http"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/api"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
)

// webDir определяет путь к директории со статическими файлами фронтенда (index.html, js, css)
const webDir = "./web"

// New создаёт и инициализирует новый маршрутизатор http.ServeMux.
// Принимает указатель на экземпляр базы данных db (*database.DB),
// регистрирует все HTTP-обработчики приложения (через функцию Init)
// и возвращает сконфигурированный ServeMux, готовый к передаче в HTTP-сервер.
func New(db *database.DB) *http.ServeMux {
	mux := http.NewServeMux()
	Init(mux, db)
	return mux
}

// Init регистрирует API-маршруты в переданном роутере mux.
// Создаёт экземпляр api.Handler, связывая его с подключением к БД (*database.DB).
func Init(mux *http.ServeMux, db *database.DB) {
	handler := api.NewHandler(db)

	// Регистрация раздачи статических файлов фронтенда
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	// Регистрация API-обработчиков
	mux.HandleFunc("/api/nextdate", api.NextDayHandler)
	mux.HandleFunc("/api/task", handler.TaskHandler)
	mux.HandleFunc("/api/tasks", handler.GetTasksHandler)
	// ...
}
