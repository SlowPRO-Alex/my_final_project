package api

import (
	"net/http"
)

const DFormat = "20060102"

type EmptyStruct struct{}

func Init() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", SignInHandler)

	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
}
