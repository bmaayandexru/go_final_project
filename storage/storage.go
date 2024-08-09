package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"` // omitempty
	Title   string `json:"title"`
	Comment string `json:"comment"` // omitempty
	Repeat  string `json:"repeat"`  // omitempty
}

const (
	templ = "20060102"
	limit = 50
)

type TaskStore struct {
	DB *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{DB: db}
}

func (ts TaskStore) Add(task Task) (sql.Result, error) {
	return ts.DB.Exec("INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?) ",
		task.Date, task.Title, task.Comment, task.Repeat)
}

func (ts TaskStore) Delete(id string) (sql.Result, error) {
	return ts.DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
}

func (ts TaskStore) Find(search string) (*sql.Rows, error) {
	// возвращаем всё если строка пустая
	if len(search) == 0 {
		return ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
			sql.Named("limit", 50))
	}
	// парсим строку на дату
	if date, err := time.Parse("02-01-2006", search); err == nil {
		// дата есть
		return ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit",
			sql.Named("date", date.Format(templ)),
			sql.Named("limit", limit))
	} else {
		// даты нет
		search = "%" + search + "%"
		return ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE UPPER(title) LIKE UPPER(:search) OR UPPER(comment) LIKE UPPER(:search) ORDER BY date LIMIT :limit",
			sql.Named("search", search),
			sql.Named("limit", limit))
	}
}

func (ts TaskStore) Get(id string) (Task, error) {
	row := ts.DB.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	task := Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, err
}

func (ts TaskStore) Update(task Task) (sql.Result, error) {
	return ts.DB.Exec("UPDATE scheduler SET  date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
}

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

var DBFileRun = "scheduler.db"

func InitDBase() (*sql.DB, error) {
	var SqlDB *sql.DB
	var StrDBFile string
	fmt.Println("Init Data Base...")
	envDBFile := os.Getenv("TODO_DBFILE")
	if envDBFile == "" {
		envDBFile = DBFileRun
		//envDBFile = tests.DBFile
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
	SqlDB.SetMaxIdleConns(5) //2
	SqlDB.SetMaxOpenConns(5)
	SqlDB.SetConnMaxIdleTime(time.Minute * 5)
	SqlDB.SetConnMaxLifetime(time.Hour)
	return SqlDB, err
}
