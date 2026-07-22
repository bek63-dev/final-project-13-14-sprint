// Package database: файл task.go содержит структуру задачи для слоя БД
// и методы для работы с таблицей scheduler
package database

import (
	"database/sql"
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
