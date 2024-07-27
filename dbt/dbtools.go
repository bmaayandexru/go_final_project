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
	ID      string `json:"id"`
	Date    string `json:"date"` // omitempty
	Title   string `json:"title"`
	Comment string `json:"comment"` // omitempty
	Repeat  string `json:"repeat"`  // omitempty
}

var SqlDB *sql.DB

func InitDBase() {
	fmt.Println("Init Data Base...")
	envDBFile := os.Getenv("TODO_DBFILE")
	if envDBFile == "" {
		envDBFile = tests.DBFileRun
	}
	fmt.Println("Result DBFile ", envDBFile)
	_, err := os.Stat(envDBFile)
	install := (err != nil)
	fmt.Println("Need install ", install)
	SqlDB, err = sql.Open("sqlite", envDBFile)
	if err != nil {
		fmt.Println("InitDB err:", err)
		return
	}
	if install {
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

func UpdateTask(task Task) (sql.Result, error) {
	return SqlDB.Exec("UPDATE scheduler SET  date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
}

func SelectID(id string) error {
	row := SqlDB.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	task := Task{}
	return row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
}
