// Package api: файл validate.go содержит функции-хелперы для валидации http-запросов
package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/domain"
)

// checkDate проверяет и при необходимости корректирует task.Date
func checkDate(task *database.Task) error {
	now := domain.TruncateDate(time.Now())

	// Если дата пустая - подставляется сегодняшняя дата
	if strings.TrimSpace(task.Date) == "" {
		task.Date = now.Format(DateLayout)
	}

	// Проверка формата даты
	t, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("time.Parse: некорректный формат исходной даты %q, ожидаемый формат %q",
			task.Date, DateLayout)
	}

	// Если задано task.Repeat — валидируется через domain.NextDate,
	// заодно вычисляется дата следующего повторения
	var next string
	if strings.TrimSpace(task.Repeat) != "" {
		next, err = domain.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %w", err)
		}
	}

	// Если дата оказалась в прошлом — заменяется на сегодняшний день (без правила повторения),
	// либо на вычисленную следующую дату (с правилом повторения)
	if domain.AfterNow(now, t) {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = now.Format(DateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}

// isMethodAllowed проверяет, соответствует ли метод запроса ожидаемому.
// Если нет — отправляет ошибку 405 Method Not Allowed и возвращает false.
func isMethodAllowed(expectedMethod string, res http.ResponseWriter, req *http.Request) bool {
	if req.Method != expectedMethod {
		writeError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
		return false
	}
	return true
}
