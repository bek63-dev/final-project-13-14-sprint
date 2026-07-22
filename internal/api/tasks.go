package api

import (
	"log"
	"net/http"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
)

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

const maxTasks = 50

// GetTasksHandler обрабатывает GET /api/tasks: возвращает список ближайших задач
func (h *Handler) GetTasksHandler(res http.ResponseWriter, req *http.Request) {
	if !isMethodAllowed(http.MethodGet, res, req) {
		return
	}

	dbTasks, err := h.DB.Tasks(maxTasks)
	if err != nil {
		log.Println("GetTasksHandler: h.DB.GetTasks: ", err)
		writeError(res, http.StatusInternalServerError, "ошибка получения списка задач")
		return
	}

	tasks := make([]Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, convToAPITask(dbTask))
	}

	writeJSON(res, http.StatusOK, TasksResponse{Tasks: tasks})

}

// convToAPITask преобразует database.Task в api.Task для отдачи через JSON
func convToAPITask(task *database.Task) Task {
	return Task{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
