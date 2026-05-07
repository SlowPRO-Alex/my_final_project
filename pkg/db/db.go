package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "", 
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX task_date ON scheduler (date);`

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Println(err)
		return err
	}

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			log.Println(err)
			return err
		}
		install = false
	}
	return nil
}
