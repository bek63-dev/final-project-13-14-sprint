package api

import (
	"net/http"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
)

// Handler - контейнер зависимостей для обработчиков API (доступ к БД)
type Handler struct {
	DB *database.DB
}

// NewHandler создаёт Handler с переданным подключением к БД
func NewHandler(db *database.DB) *Handler {
	return &Handler{DB: db}
}

// TaskHandler — точка входа для "/api/task", диспетчеризация по методу
func (h *Handler) TaskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		h.AddTaskHandler(res, req)
	// GET/PUT/DELETE будут добавлены позже

	default:
		writeError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}
