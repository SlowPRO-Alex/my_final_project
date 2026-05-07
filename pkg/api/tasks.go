package api

import (
	"net/http"

	"github.com/SlowPRO-Alex/my_final_project/pkg/db"
)

const limit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	search := params.Get("search")
	tasks, err := db.Tasks(limit, search) // в параметре максимальное количество записей и строка поиска
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
