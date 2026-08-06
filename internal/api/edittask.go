// Package api: файл edittask.go содержит обработчики получения и редактирования
// одной задачи по идентификатору (GET и PUT /api/task).
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// idParam — имя GET-параметра идентификатора задачи, общее для всех
// обработчиков, принимающих id (получение, удаление, отметка выполнения).
const idParam = "id"

// GetTaskHandler обрабатывает GET /api/task?id=<id>: возвращает задачу по идентификатору.
// Используется фронтендом для заполнения формы редактирования задачи.
func (h *Handler) GetTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	id := req.FormValue(idParam)
	if strings.TrimSpace(id) == "" {
		writeError(res, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	dbTask, err := h.DB.GetTask(id)
	if err != nil {
		log.Println("GetTaskHandler: h.DB.GetTask:", err)
		writeError(res, http.StatusBadRequest, "задача не найдена")
		return
	}

	writeJSON(res, http.StatusOK, convToAPITask(dbTask))
}

// EditTaskHandler обрабатывает PUT /api/task: принимает JSON с
// идентификатором задачи, валидирует поля и обновляет запись в БД.
func (h *Handler) EditTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		writeError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	var task Task
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		log.Println("EditTaskHandler: json.NewDecoder:", err)
		writeError(res, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if strings.TrimSpace(task.ID) == "" {
		writeError(res, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(res, http.StatusBadRequest, "не указан заголовок задачи")
		return
	}

	dbTask := convToDBTask(&task)

	if err := checkDate(dbTask); err != nil {
		log.Println("EditTaskHandler: checkDate:", err)
		writeError(res, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.DB.UpdateTask(dbTask); err != nil {
		log.Println("EditTaskHandler: h.DB.UpdateTask:", err)
		writeError(res, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(res, http.StatusOK, struct{}{})
}
