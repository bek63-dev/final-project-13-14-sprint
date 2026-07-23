// Package database: файл task.go содержит структуру задачи для слоя БД
// и методы для работы с таблицей scheduler
package database

import (
	"strconv"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/domain"
)

// Task описывает задачи на уровне БД (не путать с api.Task, который принимает/отдаёт JSON через HTTP).
// Такое разделение позволяет менять формат API, не трогая структуру хранения.
type Task struct {
	ID      string
	Date    string
	Title   string
	Comment string
	Repeat  string
}

// AddTask сохраняет задачу в таблице scheduler и возвращает идентификатор новой записи (LastInsertId).
func (db *DB) AddTask(task *Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Tasks ищет задачи по дате или ключевому слову в title/comment
// и возвращает список, отсортированный по дате с ограничением limit.
func (db *DB) Tasks(search string, limit int) ([]*Task, error) {
	var query string
	var args []any

	switch {
	case search == "":
		query = "SELECT * FROM scheduler ORDER BY date LIMIT ?"
		args = []any{limit}
	default:
		if date, ok := domain.ParseSearchDate(search); ok {
			query = "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
			args = []any{date, limit}
		} else {
			search = "%" + search + "%"
			query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
			args = []any{search, search, limit}
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		var id int64
		task := &Task{}

		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
