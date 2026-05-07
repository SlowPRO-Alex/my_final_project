package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SlowPRO-Alex/my_final_project/pkg/api"
	"github.com/SlowPRO-Alex/my_final_project/pkg/db"
)

const defaultPort = "7540"

func main() {
	sPort := os.Getenv("TODO_PORT")
	if sPort == "" {
		sPort = defaultPort
	} else {
		log.Printf("TODO_PORT = %s\n", sPort)
	}
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	} else {
		log.Printf("TODO_DBFILE = %s\n", dbFile)
	}
	if os.Getenv("TODO_PASSWORD") != "" {
		log.Printf("TODO_PASSWORD = %s\n", os.Getenv("TODO_PASSWORD"))
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	log.Printf("Start server. localhost:%s\n", sPort)
	api.Init()
	err = http.ListenAndServe(fmt.Sprintf(":%s", sPort), nil)
	if err != nil {
		panic(err)
	}

}
