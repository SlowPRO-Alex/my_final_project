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
	now := time.Now().UTC()
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
	if afterNow(now, t) {
		if task.Repeat == "" {
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
	fmt.Println("New task added: ", map[string]string{"id": fmt.Sprintf("%s", strconv.FormatInt(res, 10))})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}
