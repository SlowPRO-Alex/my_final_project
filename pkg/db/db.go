package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"fmt"
	"os"
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

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	db, err := sql.Open("sqlite", dbFile)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer db.Close()
	
	if install {
	_, err := db.Exec(schema)
	if err != nil {
        fmt.Println(err)
        return err
    }
	}
	//fmt.Println(db)
	return nil
}