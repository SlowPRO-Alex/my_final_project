package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SlowPRO-Alex/my_final_project/pkg/db"
)

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
	}
}

func checkDate(task db.Task) (date string, err error) {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var next string
	if task.Date == "" {
		task.Date = now.Format(DFormat)
	}
	t, err := time.Parse(DFormat, task.Date)
	if err != nil {
		fmt.Println(err)
	}
	next, err = NextDate(now, task.Date, task.Repeat)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(next, t)
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DFormat)
		} else {
			task.Date = next
		}
	}
	return task.Date, err
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан title"})
		return
	}
	task.Date, err = checkDate(task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	res, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{"id": fmt.Sprintf("%s", strconv.FormatInt(res, 10))})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		_, err := strconv.Atoi(id)
		if id != "" && err == nil {
			err := db.DeleteTask(id)
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()})
				return
			}
			writeJson(w, EmptyStruct{})
			return
		} else {
			writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		}
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		if r.URL.Query().Get("id") != "" {
			task, err := db.GetTask(r.URL.Query().Get("id"))
			if err != nil {
				fmt.Println("Задача не найдена:", err)
				writeJson(w, map[string]string{"error": "Задача не найдена"})
				return
			}
			writeJson(w, task)
		} else {
			writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		}
	case http.MethodPut:
		var task db.Task
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		task.Date, err = checkDate(task)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		err = db.UpdateTask(&task)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, EmptyStruct{})

	}
}
