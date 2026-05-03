package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/SlowPRO-Alex/my_final_project/pkg/api"
	"github.com/SlowPRO-Alex/my_final_project/pkg/db"
	"github.com/SlowPRO-Alex/my_final_project/tests"
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

	fmt.Printf("Start server.\nlocalhost:%s\n", sPort)
	api.Init()
	err = http.ListenAndServe(fmt.Sprintf(":%s", sPort), nil)
	if err != nil {
		panic(err)
	}

}
