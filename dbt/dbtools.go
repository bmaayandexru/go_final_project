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
	"time"

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
var StrDBFile string

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
	StrDBFile = envDBFile
	SqlDB, err = sql.Open("sqlite", StrDBFile)
	if err != nil {
		fmt.Println("InitDB err:", err)
		return
	}
	if install {
		if _, err = SqlDB.Exec(schemaSQL); err != nil {
			fmt.Println("InitDB err:", err)
		}
	}
	defer SqlDB.Close()
}

func AddTask(task Task) (sql.Result, error) {
	SqlDB, _ = sql.Open("sqlite", StrDBFile)
	defer SqlDB.Close()
	return SqlDB.Exec("INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?) ",
		task.Date, task.Title, task.Comment, task.Repeat)
}

func UpdateTask(task Task) (sql.Result, error) {
	SqlDB, _ = sql.Open("sqlite", StrDBFile)
	defer SqlDB.Close()
	return SqlDB.Exec("UPDATE scheduler SET  date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
}

func SelectByID(id string) (Task, error) {
	SqlDB, _ = sql.Open("sqlite", StrDBFile)
	defer SqlDB.Close()
	row := SqlDB.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	task := Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, err
}

// запрос по дате
func QueryByDate(date time.Time) (*sql.Rows, error) {
	SqlDB, _ = sql.Open("sqlite", StrDBFile)
	defer SqlDB.Close()
	return SqlDB.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit",
		sql.Named("date", date.Format("20060102")),
		sql.Named("limit", 50))
}

func QueryByString(search string) (*sql.Rows, error) {
	SqlDB, _ = sql.Open("sqlite", StrDBFile)
	defer SqlDB.Close()
	search = "%" + search + "%"
	return SqlDB.Query("SELECT * FROM scheduler WHERE UPPER(title) LIKE UPPER(:search) OR UPPER(comment) LIKE UPPER(:search) ORDER BY date LIMIT :limit",
		sql.Named("search", search),
		sql.Named("limit", 50))
}

// запрос всех задач
func QueryAllTasks() (*sql.Rows, error) {
	SqlDB, _ = sql.Open("sqlite", StrDBFile)
	defer SqlDB.Close()
	return SqlDB.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", 50))
}

func DeleteByID(id string) (sql.Result, error) {
	SqlDB, _ = sql.Open("sqlite", StrDBFile)
	defer SqlDB.Close()
	return SqlDB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
}
