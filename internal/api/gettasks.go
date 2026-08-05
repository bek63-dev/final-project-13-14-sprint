// Package api: файл gettasks.go содержит обработчик получения списка задач
// (GET /api/tasks) с поддержкой поиска по дате или ключевому слову.
package api

import (
	"log"
	"net/http"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
)

// TasksResponse — формат ответа со списком задач: {"tasks": [...]}.
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

const (
	maxTasks     = 50       // - ограничивает число задач, возвращаемых за один запрос к /api/tasks.
	searchParams = "search" // - имя GET-параметра запроса для поиска/фильтрации задач
)

// GetTasksHandler обрабатывает GET /api/tasks - возвращает список ближайших задач (не более maxTasks).
// Принимает необязательный query-параметр search:
//   - если "ДД.ММ.ГГГГ": ищет задачи на конкретную дату;
//   - иначе: ищет по совпадению строки в названии (title) или комментарии (comment);
//   - если пустой: возвращает ближайшие задачи отсортированные по дате по возрастанию.
func (h *Handler) GetTasksHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	search := req.FormValue(searchParams)

	dbTasks, err := h.DB.Tasks(search, maxTasks)
	if err != nil {
		log.Println("GetTasksHandler: h.DB.Tasks: ", err)
		writeError(res, http.StatusInternalServerError, "ошибка получения списка задач")
		return
	}

	tasks := make([]Task, 0, len(dbTasks))

	for _, dbTask := range dbTasks {
		tasks = append(tasks, convToAPITask(dbTask))
	}

	writeJSON(res, http.StatusOK, TasksResponse{Tasks: tasks})
}

// convToAPITask преобразует database.Task в Task для передачи клиенту через JSON.
func convToAPITask(task *database.Task) Task {
	return Task{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
