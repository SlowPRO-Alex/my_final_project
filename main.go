package main

import (
    "fmt"
    "net/http"
	"github.com/SlowPRO-Alex/my_final_project/tests"
	"github.com/SlowPRO-Alex/my_final_project/pkg/db"
	"github.com/SlowPRO-Alex/my_final_project/pkg/api"
	"os"
	"strconv"
)

func main() {
	sPort := os.Getenv("TODO_PORT")
	dbFile := os.Getenv("TODO_DBFILE")
	if sPort == "" {
		sPort = strconv.Itoa(tests.Port)
	}
	if dbFile == "" {
		dbFile = tests.DBFile
	}
	
	err := db.Init(dbFile)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Запускаем сервер.\nАдрес: localhost:%s\n",sPort)
	api.Init()
    err = http.ListenAndServe(fmt.Sprintf(":%s",sPort), nil)
    if err != nil {
        panic(err)
    }

} 