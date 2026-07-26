// Package api: файл donetask.go содержит обработчики отметки выполненной задачи
// (POST /api/task/done) и удаления задачи (DELETE /api/task).
package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/domain"
)

// DoneTaskHandler обрабатывает POST /api/task/done?id=<id>: отмечает
// задачу выполненной. Одноразовая задача (пустой Repeat) удаляется,
// периодическая — переносится на следующую дату по своему правилу.
func (h *Handler) DoneTaskHandler(res http.ResponseWriter, req *http.Request) {
	if !isMethodAllowed(http.MethodPost, res, req) {
		return
	}

	id := req.FormValue(idParam)
	if strings.TrimSpace(id) == "" {
		writeError(res, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	dbTask, err := h.DB.GetTask(id)
	if err != nil {
		log.Println("DoneTaskHandler: h.DB.GetTask:", err)
		writeError(res, http.StatusBadRequest, "задача не найдена")
		return
	}

	// Одноразовая задача считается выполненной и удаляется из БД.
	if strings.TrimSpace(dbTask.Repeat) == "" {
		if err := h.DB.DeleteTask(id); err != nil {
			log.Println("DoneTaskHandler: h.DB.DeleteTask:", err)
			writeError(res, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(res, http.StatusOK, struct{}{})
		return
	}

	// Периодическая задача переносится на следующую дату по правилу Repeat.
	now := domain.TruncateDate(time.Now())

	next, err := domain.NextDate(now, dbTask.Date, dbTask.Repeat)
	if err != nil {
		log.Println("DoneTaskHandler: domain.NextDate:", err)
		writeError(res, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.DB.UpdateDate(next, id); err != nil {
		log.Println("DoneTaskHandler: h.DB.UpdateDate:", err)
		writeError(res, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(res, http.StatusOK, struct{}{})
}

// DeleteTaskHandler обрабатывает DELETE /api/task?id=<id>: удаляет задачу по идентификатору.
func (h *Handler) DeleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	if !isMethodAllowed(http.MethodDelete, res, req) {
		return
	}

	id := req.FormValue(idParam)
	if strings.TrimSpace(id) == "" {
		writeError(res, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	if err := h.DB.DeleteTask(id); err != nil {
		log.Println("DeleteTaskHandler: h.DB.DeleteTask:", err)
		writeError(res, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(res, http.StatusOK, struct{}{})
}
