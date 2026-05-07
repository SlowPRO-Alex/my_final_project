package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SlowPRO-Alex/my_final_project/pkg/db"
)

func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
	}
	next, err = NextDate(now, task.Date, task.Repeat)
	if err != nil {
		log.Println(err)
	}
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(body, &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title is not specified"}, http.StatusBadRequest)
		return
	}
	task.Date, err = checkDate(task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	res, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, map[string]string{"id": fmt.Sprintf("%s", strconv.FormatInt(res, 10))}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := strconv.Atoi(id)
	if id != "" && err == nil {
		err := db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJson(w, EmptyStruct{}, http.StatusOK)
		return
	}
	writeJson(w, map[string]string{"error": "id is not specified"}, http.StatusBadRequest)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("id") != "" {
		task, err := db.GetTask(r.URL.Query().Get("id"))
		if err != nil {
			log.Println("Task not found:", err)
			writeJson(w, map[string]string{"error": "task not found"}, http.StatusBadRequest)
			return
		}
		writeJson(w, task, http.StatusOK)
		return
	}
	writeJson(w, map[string]string{"error": "id is not specified"}, http.StatusBadRequest)

}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	task.Date, err = checkDate(task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	writeJson(w, EmptyStruct{}, http.StatusOK)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	default:
		writeJson(w, map[string]string{"error": "Method Not Allowed"}, http.StatusMethodNotAllowed)
	}
}
