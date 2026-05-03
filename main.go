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
	if sPort == "" {
		sPort = strconv.Itoa(tests.Port)
	} else {
		fmt.Printf("TODO_PORT = %s\n", sPort)
	}
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = tests.DBFile
	} else {
		fmt.Printf("TODO_DBFILE = %s\n", dbFile)
	}
	if os.Getenv("TODO_PASSWORD") != "" {
		fmt.Printf("TODO_PASSWORD = %s\n", os.Getenv("TODO_PASSWORD"))
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
