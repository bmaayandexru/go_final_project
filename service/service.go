package service

import (
	"database/sql"
	"time"

	"github.com/bmaayandexru/go_final_project/dbt"
)

type TaskStore struct {
	DB *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{DB: db}
}

type TaskService struct {
	store TaskStore
}

const limit = 50

var Service TaskService

func InitStoreAndService(db *sql.DB) {
	Service = NewTaskService(NewTaskStore(db))
}

func NewTaskService(store TaskStore) TaskService {
	return TaskService{store: store}
}

func (ts TaskStore) Add(task dbt.Task) (sql.Result, error) {
	return ts.DB.Exec("INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?) ",
		task.Date, task.Title, task.Comment, task.Repeat)
}

func (ts TaskStore) Delete(id string) (sql.Result, error) {
	return ts.DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
}

func (ts TaskStore) Find(search string) (*sql.Rows, error) {
	// парсим строку на дату
	if date, err := time.Parse("02-01-2006", search); err == nil {
		// дата есть
		return ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit",
			sql.Named("date", date.Format("20060102")),
			sql.Named("limit", limit))
	} else {
		// даты нет
		search = "%" + search + "%"
		return ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE UPPER(title) LIKE UPPER(:search) OR UPPER(comment) LIKE UPPER(:search) ORDER BY date LIMIT :limit",
			sql.Named("search", search),
			sql.Named("limit", limit))
	}
}

func (ts TaskStore) QueryAllTasks() (*sql.Rows, error) {
	return ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
		sql.Named("limit", 50))
}

func (ts TaskStore) Select(id string) (dbt.Task, error) {
	row := ts.DB.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	task := dbt.Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, err
}

func (ts TaskStore) Update(task dbt.Task) (sql.Result, error) {
	return ts.DB.Exec("UPDATE scheduler SET  date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
}

func (ts TaskService) Add(task dbt.Task) (sql.Result, error) {
	return ts.store.Add(task)
}

func (ts TaskService) Delete(id string) (sql.Result, error) {
	return ts.store.Delete(id)
}

func (ts TaskService) Find(search string) (*sql.Rows, error) {
	return ts.store.Find(search)
}

func (ts TaskService) QueryAllTasks() (*sql.Rows, error) {
	return ts.store.QueryAllTasks()
}

func (ts TaskService) Select(id string) (dbt.Task, error) {
	return ts.store.Select(id)
}

func (ts TaskService) Update(task dbt.Task) (sql.Result, error) {
	return ts.store.Update(task)
}
