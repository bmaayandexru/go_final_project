package dbt

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var schemaSQL string = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "20000101", 
    title VARCHAR(32) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX idx_date ON scheduler (date); 
CREATE INDEX idx_title ON scheduler (title); 
`

// CREATE UNIQUE INDEX idx_title ON scheduler (title);

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"` // omitempty
	Title   string `json:"title"`
	Comment string `json:"comment"` // omitempty
	Repeat  string `json:"repeat"`  // omitempty
}

const limit = 50
const template = "20060102"

// путь к базе для запуска
// эту переменную использую у себя в коде
var DBFileRun = "scheduler.db"

var SqlDB *sql.DB
var StrDBFile string

func InitDBase() (*sql.DB, error) {
	fmt.Println("Init Data Base...")
	envDBFile := os.Getenv("TODO_DBFILE")
	if envDBFile == "" {
		envDBFile = DBFileRun
	}
	fmt.Println("Result DBFile ", envDBFile)
	_, err := os.Stat(envDBFile)
	install := (err != nil)
	fmt.Println("Need install ", install)
	StrDBFile = envDBFile
	SqlDB, err = sql.Open("sqlite", StrDBFile)
	if err != nil {
		fmt.Println("InitDB err:", err)
		return SqlDB, err
	}
	if install {
		if _, err = SqlDB.Exec(schemaSQL); err != nil {
			fmt.Println("InitDB err:", err)
			// SqlDB = nil
			return SqlDB, err
		}
	}
	return SqlDB, err
}
