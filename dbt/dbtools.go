/*
	Пакет - каталог

во всех файлах *.go внутри пакета-каталога объявляем pakage <имя каталога>
package "dbt"
имена файлов не важны
импортируются пакеты import "github.com/<github user name>/<project name>/<packege name>"
import "github.com/bmaayandexru/go_final_project/dbt"
пускаем обязательно из модуля main
изменения подхватятся автоматически
заливать предварительно на гитхаб ничего не надо
*/
package dbt

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	"github.com/bmaayandexru/go_final_project/tests"
)

// см bmanote.md
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
	ID      int    `json:"id"`
	Date    string `json:"date"` // omitempty
	Title   string `json:"title"`
	Comment string `json:"comment"` // omitempty
	Repeat  string `json:"repeat"`  // omitempty
}

var SqlDB *sql.DB

func InitDBase() {
	fmt.Println("Init Data Base...")
	//fmt.Println("settings DBFile", tests.DBFile)
	envDBFile := os.Getenv("TODO_DBFILE")
	//fmt.Println("enviroment DBFile ", envDBFile)
	if envDBFile == "" {
		envDBFile = tests.DBFileRun
	}
	fmt.Println("Result DBFile ", envDBFile)
	_, err := os.Stat(envDBFile)
	//	var install bool
	install := (err != nil)
	fmt.Println("Need install ", install)
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	SqlDB, err = sql.Open("sqlite", envDBFile)
	if err != nil {
		fmt.Println("InitDB err:", err)
		return
	}
	if install {
		// нужно создать таблицу, т к файла не было
		if _, err = SqlDB.Exec(schemaSQL); err != nil {
			fmt.Println("InitDB err:", err)
		}
	}

	// defer sqlDB.Close()
}

func AddTask(task Task) (sql.Result, error) {
	return SqlDB.Exec("INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?) ",
		task.Date, task.Title, task.Comment, task.Repeat)
}
