// Package database: файл task.go содержит структуру задачи для слоя БД
// и методы для работы с таблицей scheduler
package database

import (
	"database/sql"
	"strconv"
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
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) Tasks(limit int) ([]*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit"

	rows, err := db.Query(query, sql.Named("limit", limit))
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
