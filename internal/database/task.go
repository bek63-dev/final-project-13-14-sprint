// Package database: файл task.go содержит модель задачи уровня БД и методы доступа к таблице scheduler.
package database

import (
	"fmt"
	"strconv"

	"github.com/bek63-dev/final-project-13-14-sprint/internal/domain"
)

// Task описывает задачу на уровне БД (не путать с api.Task, который
// принимает и отдаёт данные через JSON). Разделение позволяет менять
// формат API, не затрагивая структуру хранения.
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

// Tasks возвращает список задач, отсортированный по дате, с ограничением limit.
// - Пустой search возвращает ближайшие задачи;
// - search в формате даты фильтрует по конкретному дню, иначе — по ключевым словам в title или comment.
func (db *DB) Tasks(search string, limit int) ([]*Task, error) {
	const selectFields = "SELECT id, date, title, comment, repeat FROM scheduler"

	date, isDate := domain.ParseSearchDate(search)

	var query string
	var args []any

	switch {
	case search == "":
		query = selectFields + " ORDER BY date LIMIT ?"
		args = []any{limit}
	case isDate:
		query = selectFields + " WHERE date == ? ORDER BY date LIMIT ?"
		args = []any{date, limit}
	default:
		like := "%" + search + "%"
		query = selectFields + " WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		args = []any{like, like, limit}
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

// GetTask возвращает задачу по её идентификатору.
func (db *DB) GetTask(id string) (*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id == ?"

	var rowID int64
	task := &Task{}

	err := db.QueryRow(query, id).Scan(&rowID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("GetTask: %w", err)
	}

	task.ID = strconv.FormatInt(rowID, 10)
	return task, nil
}

// UpdateTask обновляет задачу по её идентификатору (task.ID).
// Возвращает ошибку, если задачи с таким id не существует.
func (db *DB) UpdateTask(task *Task) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id == ?"

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("UpdateTask: db.Exec: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateTask: RowsAffected: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("некорректный идентификатор для обновления задачи")
	}

	return nil
}

// DeleteTask удаляет задачу по переданному идентификатору.
// Возвращает ошибку, если задачи с таким id не существует.
func (db *DB) DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id == ?"

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("DeleteTask: db.Exec: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteTask: RowsAffected: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("некорректный идентификатор для удаления задачи")
	}

	return nil
}

// UpdateDate обновляет только дату задачи с указанным идентификатором, не
// затрагивая остальные поля. Используется при выполнении периодической
// задачи, когда дату нужно сдвинуть на следующую по правилу повторения.
func (db *DB) UpdateDate(next string, id string) error {
	query := "UPDATE scheduler SET date = ? WHERE id == ?"

	res, err := db.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("UpdateDate: db.Exec: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateDate: RowsAffected: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("некорректный идентификатор для обновления даты задачи")
	}

	return nil
}
