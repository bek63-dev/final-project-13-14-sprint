// Package api: файл addtask.go содержит обработчик для добавления новой задачи (POST /api/task)
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
)

// Task описывает задачу для приёма и передачи данных через JSON API.
// Отделён от database.Task, чтобы формат API можно было менять (JSON-теги, набор полей),
// не трогая структуру хранения в БД.
type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// AddTaskHandler обрабатывает POST /api/task: принимает JSON, валидирует его,
// корректирует дату выполнения и сохраняет задачу в БД.
// Ответ:
//
//	{"id": "123"}       — при успехе
//	{"error": "текст"}  — при любой ошибке
func (h *Handler) AddTaskHandler(res http.ResponseWriter, req *http.Request) {
	if !isMethodAllowed(http.MethodPost, res, req) {
		return
	}

	var task Task
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		log.Println("AddTaskHandler: json.NewDecoder: ", err)
		writeError(res, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}
	defer req.Body.Close()

	// Валидация обязательного поля Title
	if strings.TrimSpace(task.Title) == "" {
		writeError(res, http.StatusBadRequest, "не указан заголовок задачи")
		return
	}

	dbTask := convToDBTask(&task)

	if err := checkDate(dbTask); err != nil {
		log.Println("AddTaskHandler: checkDate:", err)
		writeError(res, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.DB.AddTask(dbTask)
	if err != nil {
		log.Println("AddTaskHandler: h.DB.AddTask:", err)
		writeError(res, http.StatusInternalServerError, "ошибка сохранения задачи")
		return
	}

	writeJSON(res, http.StatusOK, idResponse{ID: strconv.FormatInt(id, 10)})
}

// convToDBTask преобразует Task в database.Task для сохранения в БД.
func convToDBTask(task *Task) *database.Task {
	return &database.Task{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
