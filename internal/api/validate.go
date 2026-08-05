// Package api: файл validate.go содержит общие утилиты валидации входящих HTTP-запросов,
// используемые обработчиками добавления и редактирования задач.
package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/database"
	"github.com/bek63-dev/final-project-13-14-sprint/internal/domain"
)

// checkDate проверяет и при необходимости корректирует дату задачи task:
// подставляет сегодняшнюю дату при пустом значении, а для дат из прошлого
// вычисляет ближайшую подходящую дату по правилу повторения.
// Вызывается перед сохранением задачи в БД при добавлении и редактировании.
func checkDate(task *database.Task) error {
	now := domain.TruncateDate(time.Now())

	if strings.TrimSpace(task.Date) == "" {
		task.Date = now.Format(DateLayout)
	}

	t, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("time.Parse: некорректный формат исходной даты %q, ожидаемый формат %q",
			task.Date, DateLayout)
	}

	// Если указано правило повторения, заранее вычисляем дату следующего
	// срабатывания — она понадобится, если исходная дата окажется в прошлом.
	var next string
	if strings.TrimSpace(task.Repeat) != "" {
		next, err = domain.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %w", err)
		}
	}

	// Дата в прошлом заменяется на сегодняшнюю (одноразовая задача) либо
	// на вычисленную следующую дату (задача с повторением).
	if domain.AfterNow(now, t) {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = now.Format(DateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}
