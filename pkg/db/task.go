package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
		return id, nil
	}
	return 0, err
}

func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
	tasks := make([]*Task, 0)

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}

		var idInt64 int64
		var dateStr, titleStr, commentStr, repeatStr string

		err := rows.Scan(&idInt64, &dateStr, &titleStr, &commentStr, &repeatStr)
		if err != nil {
			return nil, err
		}

		task.ID = strconv.FormatInt(idInt64, 10)
		task.Date = dateStr
		task.Title = titleStr
		task.Comment = commentStr
		task.Repeat = repeatStr

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	log.Printf("[DB DEBUG] Возвращаем слайс задач. Длина: %d, значение: %v", len(tasks), tasks)
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := "SELECT date, title, comment, repeat FROM scheduler WHERE id = ?"

	var task Task
	task.ID = id

	row := db.QueryRow(query, id)

	err := row.Scan(&task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := ` UPDATE scheduler SET title = ?, comment = ?, date = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Title, task.Comment, task.Date, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := ` DELETE FROM scheduler WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := ` UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
