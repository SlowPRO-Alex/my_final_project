package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const DFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, req *http.Request) {
	now := time.Now().UTC()
	sNow := req.URL.Query().Get("now")
	if len(sNow) > 0 {
		pNow, err := time.Parse(DFormat, sNow)
		now = pNow

		if err != nil {
			fmt.Println(err)
		}
	}
	dstart := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")
	answer, err := NextDate(now, dstart, repeat)
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(w, answer)
}

func Init() {
	http.Handle(`/`, http.FileServer(http.Dir("./web/")))
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
}
