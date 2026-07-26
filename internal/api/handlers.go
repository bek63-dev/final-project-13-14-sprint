// Package api реализует транспортный слой планировщика задач:
// HTTP-обработчики аутентификации и CRUD-операций над задачами.
// Обработчики принимают запросы, валидируют вход и делегируют работу пакетам database и domain.
package api

import (
	"net/http"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/config"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
)

// Handler хранит зависимости, общие для всех обработчиков API:
// подключение к БД и конфигурацию приложения (в том числе пароль аутентификации).
type Handler struct {
	DB     *database.DB
	Config *config.Config
}

// NewHandler создаёт Handler с указанными подключением к БД и конфигурацией.
// Используется при инициализации маршрутизатора.
func NewHandler(db *database.DB, cfg *config.Config) *Handler {
	return &Handler{
		DB:     db,
		Config: cfg,
	}
}

// TaskHandler маршрутизирует запросы "/api/task" по HTTP-методу к
// соответствующему CRUD-обработчику задачи.
func (h *Handler) TaskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.AddTaskHandler(res, req)
	case http.MethodGet:
		h.GetTaskHandler(res, req)
	case http.MethodPut:
		h.EditTaskHandler(res, req)
	case http.MethodDelete:
		h.DeleteTaskHandler(res, req)
	default:
		writeError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}
