package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func DeleteTask(id string) error {
	_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	fmt.Printf("Задача %s удалена\n", id)
	return err
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query, sql.Named("date", next), sql.Named("id", id))
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

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	if task.Title == "" {
		return fmt.Errorf("Не указан title")
	}
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM scheduler WHERE id LIKE '%s'", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return &task, err
}

func AddTask(task *Task) (int64, error) {
	var id int64
	if task.Title == "" {
		return id, errors.New("Не указан заголовок задачи")
	}

	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	fmt.Println("Add new task: id", id, task)
	return id, err
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, limit)
	var query string
	if search == "" {
		query = fmt.Sprintf("SELECT * FROM scheduler ORDER BY date LIMIT %d", limit)
	} else {
		parsedDate, err := time.Parse("02.01.2006", search)
		if err != nil {
			fmt.Println("Ошибка при парсинге даты:", err)
			query = fmt.Sprintf("SELECT * FROM scheduler WHERE title LIKE '%%%s%%' OR comment LIKE '%%%s%%' ORDER BY date LIMIT %d", search, search, limit)
		} else {
			formated := parsedDate.Format("20060102")
			query = fmt.Sprintf("SELECT * FROM scheduler WHERE date LIKE '%s' ORDER BY date LIMIT %d", formated, limit)
		}
	}
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var i int
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			return tasks[0:i], err
		}
		tasks[i] = &task
		if i == limit-1 {
			return tasks[0 : i+1], err
		} else {
			i++
		}
	}
	return tasks[0:i], nil
}
