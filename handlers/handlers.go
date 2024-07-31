package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	dbt "github.com/bmaayandexru/go_final_project/dbt"
	"github.com/bmaayandexru/go_final_project/dbt/nextdate"
	nd "github.com/bmaayandexru/go_final_project/dbt/nextdate"
)

// для формирования ошибки
type strcErr struct {
	Error string `json:"error"`
}

// для формирования идентификатора
type strcId struct {
	Id string `json:"id"`
}

// для формирования строки
type strcPwd struct {
	Password string `json:"password"`
}

const (
	CPassword = "12111"
)

type strcToken struct {
	Token string `json:"token"`
}

// для формирования "пустышки"
type strcEmpty struct{}

// для формирования слайса задач
type strcTasks struct {
	Tasks []dbt.Task `json:"tasks"`
}

var sTasks strcTasks

func NextDateHandle(res http.ResponseWriter, req *http.Request) {
	strNow := req.FormValue("now")
	strDate := req.FormValue("date")
	strRepeat := req.FormValue("repeat")
	now, err := time.Parse("20060102", strNow)
	if err != nil {
		return
	}
	retStr, err := nd.NextDate(now, strDate, strRepeat)
	if err != nil {
		return
	}
	res.Write([]byte(retStr))
}

func retError(res http.ResponseWriter, sErr string, statusCode int) {
	// переименовать в retError
	var bE strcErr
	bE.Error = sErr
	aBytes, _ := json.Marshal(bE)
	// *** лог контроль
	fmt.Println("retError: aBytes ", string(aBytes))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(aBytes)
}

func TasksGETSearchString(res http.ResponseWriter, req *http.Request) {
	search := req.URL.Query().Get("search")
	fmt.Printf("Строка *%s*\n", search)
	rows, err := dbt.QueryByString(search)
	if err != nil {
		fmt.Println(err)
		retError(res, fmt.Sprintf("Ts GET SS: Ошибка запроса: %s\n", err.Error()), http.StatusOK)
		return
	}
	defer rows.Close()
	sTasks.Tasks = make([]dbt.Task, 0)
	for rows.Next() {
		task := dbt.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			retError(res, fmt.Sprintf("Ts GET SS: Ошибка rows.Scan(): %s\n", err.Error()), http.StatusOK)
			return
		}
		sTasks.Tasks = append(sTasks.Tasks, task)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		retError(res, fmt.Sprintf("Ts GET SS: Ошибка rows.Next(): %s\n", err.Error()), http.StatusOK)
		return
	}
	arrBytes, err := json.Marshal(sTasks)
	if err != nil {
		retError(res, fmt.Sprintf("TH GET SS: Ошибка json.Marshal(sTsks): %v\n", err), http.StatusOK)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(arrBytes)
}

func TasksGETSearchDate(res http.ResponseWriter, req *http.Request) {
	search := req.URL.Query().Get("search")
	date, _ := time.Parse("02.01.2006", search)
	fmt.Printf("Дата %v\n", date)
	rows, err := dbt.QueryByDate(date)
	if err != nil {
		fmt.Println(err)
		retError(res, fmt.Sprintf("Ts GET SD: Ошибка запроса: %s\n", err.Error()), http.StatusOK)
		return
	}
	defer rows.Close()
	sTasks.Tasks = make([]dbt.Task, 0)
	for rows.Next() {
		task := dbt.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			retError(res, fmt.Sprintf("Ts GET SD: Ошибка rows.Scan(): %s\n", err.Error()), http.StatusOK)
			return
		}
		sTasks.Tasks = append(sTasks.Tasks, task)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		retError(res, fmt.Sprintf("Ts GET SD: Ошибка rows.Next(): %s\n", err.Error()), http.StatusOK)
		return
	}
	arrBytes, err := json.Marshal(sTasks)
	if err != nil {
		retError(res, fmt.Sprintf("TH GET SD: Ошибка json.Marshal(sTsks): %v\n", err), http.StatusOK)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(arrBytes)
}

func TasksGETAllTasks(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("Вывести все задачи")
	rows, err := dbt.QueryAllTasks() // *** DBT ADD ***
	//dbt.SqlDB, _ = sql.Open("sqlite", dbt.StrDBFile)
	//defer dbt.SqlDB.Close()
	//rows, err := dbt.SqlDB.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", 50))
	if err != nil {
		fmt.Println(err)
		retError(res, fmt.Sprintf("Ts GET AT: Ошибка запроса: %s\n", err.Error()), http.StatusOK)
		return
	}
	defer rows.Close()
	sTasks.Tasks = make([]dbt.Task, 0)
	for rows.Next() {
		task := dbt.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			retError(res, fmt.Sprintf("Ts GET AT: Ошибка rows.Scan(): %s\n", err.Error()), http.StatusOK)
			return
		}
		sTasks.Tasks = append(sTasks.Tasks, task)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		retError(res, fmt.Sprintf("Ts GET AT: Ошибка rows.Next(): %s\n", err.Error()), http.StatusOK)
		return
	}
	arrBytes, err := json.Marshal(sTasks)
	if err != nil {
		retError(res, fmt.Sprintf("TH GET AT: Ошибка json.Marshal(sTasks): %v\n", err), http.StatusOK)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(arrBytes)
}

func TaskGETHandle(res http.ResponseWriter, req *http.Request) {
	// получаем значение GET-параметра с именем id
	id := req.URL.Query().Get("id")
	fmt.Printf("Tk GET id %s\n", id)
	task, err := dbt.SelectByID(id)
	if err != nil {
		fmt.Println(err)
		retError(res, fmt.Sprintf("Tk GET id: Ошибка row.Scan(): %s\n", err.Error()), http.StatusOK)
		return
	}
	fmt.Println("Считана задача: ", task)
	arrBytes, err := json.Marshal(task)
	if err != nil {
		retError(res, fmt.Sprintf("Tk GET id: Ошибка json.Marshal(sTsks): %v\n", err), http.StatusOK)
		return
	}
	fmt.Printf("Tk GET id ret json *%s*\n", string(arrBytes))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(arrBytes)
}

func TasksGETHandle(res http.ResponseWriter, req *http.Request) {
	// получаем значение GET-параметра с именем search
	search := req.URL.Query().Get("search")
	if len(search) == 0 { // вывести все задачи
		TasksGETAllTasks(res, req)
		return
	} else {
		fmt.Printf("Search *%s*\n", search)
		_, e := time.Parse("02.01.2006", search)
		if e == nil { // это дата
			TasksGETSearchDate(res, req)
			return
		} else { // не получилось - ищем строку в title лил comment
			TasksGETSearchString(res, req)
			return
		}
	}
}

func TaskPUTHandle(res http.ResponseWriter, req *http.Request) {
	var task dbt.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		retError(res, fmt.Sprintf("Tk PUT: Ошибка чтения тела запроса: %s\n", err.Error()), http.StatusOK)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		retError(res, fmt.Sprintf("Tk PUT: Ошибка десериализации: %s\n", err.Error()), http.StatusOK)
		return
	}
	// лог контроль
	fmt.Printf("Tk PUT Unmarshal Task: ID *%s* Date *%s* Title *%s* Comment *%s* Repeat *%s*\n", task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	// анализ task.ID не пустой, это число, есть в базе,
	if len(task.ID) == 0 { // пустой id
		retError(res, fmt.Sprintln("Tk PUT: Пустой ID"), http.StatusOK)
		return
	}
	_, err = strconv.Atoi(task.ID)
	if err != nil { // id не число
		retError(res, fmt.Sprintln("Tk PUT: ID не число"), http.StatusOK)
		return
	}
	// тут ID строка
	_, err = dbt.SelectByID(task.ID)
	if err != nil { // запрос SELECT * WHERE id = :id не должен вернуть ошибку
		retError(res, fmt.Sprintf("Tk PUT: ID нет в базе. Ошибка: %v\n", err), http.StatusOK)
		return
	}
	// ID корректный и в базе есть
	if len(task.Title) == 0 { // Поле Title обязательное
		retError(res, "Tk PUT: Поле `Задача*` пустое", http.StatusOK)
		return
	}
	if len(task.Date) == 0 { // Если поле date не указано или содержит пустую строку,
		task.Date = time.Now().Format("20060102") // берётся сегодняшнее число.
	} else {
		//  task.Date не пустое. пробуем распарсить
		_, err = time.Parse("20060102", task.Date)
		if err != nil {
			retError(res, fmt.Sprintf("Tk PUT: Ошибка разбора даты: %v\n", err), http.StatusOK)
			return
		}
	}
	// тут валидная строка в task.Date
	if task.Date < time.Now().Format("20060102") {
		if len(task.Repeat) == 0 {
			task.Date = time.Now().Format("20060102")
		} else {
			task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				retError(res, fmt.Sprintf("Tk PUT: Ошибка NextDate: %v", err), http.StatusOK)
				return
			}
		}
	} else {
		task.Date, err = nd.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			retError(res, fmt.Sprintf("Tk PUT: Ошибка NextDate: %v", err), http.StatusOK)
			return
		}
	}
	// Task перезаписать в базе
	_, err = dbt.UpdateTask(task)
	if err != nil {
		retError(res, fmt.Sprintf("Tк PUT: Ошибка при изменении в БД: %v\n", err), http.StatusOK)
		return
	}
	fmt.Println("Изменена в базе ", task)
	// всё отлично. возвращем пустышку
	var bE strcEmpty
	arrBytes, err := json.Marshal(bE)
	if err != nil {
		retError(res, fmt.Sprintf("Tk PUT: Ошибка json.Marshal(bE): %v\n", err), http.StatusOK)
		return
	}
	// *** лог контроль
	fmt.Printf("Tk PUT: ret json *%s*\n", string(arrBytes))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(arrBytes)
}

func TaskPOSTHandle(res http.ResponseWriter, req *http.Request) {
	// добавление задачи
	var task dbt.Task
	var buf bytes.Buffer
	var err error
	var bId strcId
	// читаем тело запроса
	if _, err := buf.ReadFrom(req.Body); err != nil {
		retError(res, fmt.Sprintf("Ts POST: Ошибка чтения тела запроса: %s\n", err.Error()), http.StatusOK)
		return
	}
	// десериализуем JSON в task
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		retError(res, fmt.Sprintf("Tr POST: Ошибка десериализации: %s\n", err.Error()), http.StatusOK)
		return
	}
	// лог контроль
	fmt.Printf("Tr POST: Unmarshal Task Date *%s* Title *%s* Comment *%s* Repeat *%s*\n", task.Date, task.Title, task.Comment, task.Repeat)
	if len(task.Title) == 0 { // Поле Title обязательное
		retError(res, "Ts POST: Поле `Задача*` пустое", http.StatusOK)
		return
	}
	// Если поле date содержит пустую строку,
	if len(task.Date) == 0 { // берётся сегодняшнее число.
		task.Date = time.Now().Format("20060102")
	} else { //  task.Date не пустое. пробуем распарсить
		if _, err := time.Parse("20060102", task.Date); err != nil {
			retError(res, fmt.Sprintf("Ts POST: Ошибка разбора даты: %v\n", err), http.StatusOK)
			return
		}
	}
	// тут валидная строка в task.Date
	// это либо строка из текущей даты либо корректная строка
	nows := time.Now().Format("20060102")
	if len(task.Repeat) > 0 {
		if task.Date < nows { // правило есть и дата меньше сегодняшней
			tn, _ := time.Parse("20060102", nows)
			if task.Date, err = nextdate.NextDate(tn, task.Date, task.Repeat); err != nil {
				retError(res, fmt.Sprintf("Ts POST: Ошибка NextDate: %v", err), http.StatusOK)
				return
			}
		}
	} else { // правила повторения нет
		if task.Date < nows { // дата меньше сегодняшней
			task.Date = nows
		}
	}
	fmt.Println("Ts POST: задача добавлена в базу ", task)
	resSql, err := dbt.AddTask(task)
	if err != nil {
		retError(res, fmt.Sprintf("Ts POST: Ошибка при добавлении в БД: %v\n", err), http.StatusOK)
		return
	}
	id, err := resSql.LastInsertId()
	if err != nil {
		retError(res, fmt.Sprintf("Ts POST: Ошибка LastInsetId(): %v\n", err), http.StatusOK)
		return
	}
	bId.Id = strconv.Itoa(int(id))
	arrBytes, err := json.Marshal(bId)
	if err != nil {
		retError(res, fmt.Sprintf("Ts POST: Ошибка json.Marshal(id): %v\n", err), http.StatusOK)
		return
	}
	// *** лог контроль
	fmt.Printf("Ts POST:ret json *%s*\n", string(arrBytes))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(arrBytes)
}

func TaskDELETEHandle(res http.ResponseWriter, req *http.Request) {
	fmt.Println("запрос task DELETE")
	// получить id
	id := req.URL.Query().Get("id")
	fmt.Printf("Tk DELETE id %s\n", id)
	if len(id) == 0 {
		// нет id
		retError(res, "Tk DELETE. Нет id", http.StatusOK)
		return
	}
	if _, err := strconv.Atoi(id); err != nil {
		retError(res, "Tk DELETE. id не число", http.StatusOK)
		return
	}
	if _, err := dbt.SelectByID(id); err != nil {
		retError(res, fmt.Sprintf("Tk DELETE: id нет в базе. %s", err.Error()), http.StatusOK)
		return
	}
	// удалить по id
	_, err := dbt.DeleteByID(id) // *** DBT ADD ***
	//dbt.SqlDB, _ = sql.Open("sqlite", dbt.StrDBFile)
	//defer dbt.SqlDB.Close()
	//_, err := dbt.SqlDB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		retError(res, fmt.Sprintf("Tk DELETE: id: Ошибка удаления из базы: %s\n", err.Error()), http.StatusOK)
		return
	}
	var bE strcEmpty
	arrBytes, err := json.Marshal(bE)
	if err != nil {
		retError(res, fmt.Sprintf("Tk DELETE: Ошибка json.Marshal(): %v\n", err), http.StatusOK)
		return
	}
	// *** лог контроль
	fmt.Printf("Tk DELETE: ret json *%s*\n", string(arrBytes))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(arrBytes)
}

func TaskHandle(res http.ResponseWriter, req *http.Request) {
	// одна задача
	// запрос в строку
	s := fmt.Sprintf("Получен запрос H: %s Path: %s M: %s", req.Host, req.URL.Path, req.Method)
	fmt.Println(s)
	switch req.Method {
	case "POST": // добавление задачи
		TaskPOSTHandle(res, req)
	case "GET": // запрос для редактирование
		TaskGETHandle(res, req)
	case "PUT": // запрос на изменение
		TaskPUTHandle(res, req)
	case "DELETE": // удаление задачи
		TaskDELETEHandle(res, req)
	}
}

func TaskDonePOSTHandle(res http.ResponseWriter, req *http.Request) {
	var bE strcEmpty
	// задача выполнена
	fmt.Println("запрос task/done POST задача выполнена")
	id := req.URL.Query().Get("id")
	fmt.Printf("Tkd POST id %s\n", id)
	task, err := dbt.SelectByID(id)
	if err != nil {
		fmt.Println(err)
		retError(res, fmt.Sprintf("Tkd POST id: Ошибка row.Scan(): %s\n", err.Error()), http.StatusOK)
		return
	}
	fmt.Println("Считана задача: ", task)
	if len(task.Repeat) == 0 {
		_, err = dbt.DeleteByID(task.ID) // *** DBT ADD ***
		if err != nil {
			retError(res, fmt.Sprintf("Tkd POST id: Ошибка удаления из базы: %s\n", err.Error()), http.StatusOK)
			return
		}
	} else { // при наличии правила повторения переназначение даты и UPDATE
		dnow := time.Now()
		dnow = dnow.AddDate(0, 0, 1)
		if dnow.Format("20060102") < task.Date {
			dnow, _ = time.Parse("20060102", task.Date)
			dnow = dnow.AddDate(0, 0, 1)
		}
		newDate, err := nd.NextDate(dnow, task.Date, task.Repeat)
		if err != nil {
			retError(res, fmt.Sprintf("Tkd POST: Ошибка NextDate(): %s\n", err.Error()), http.StatusOK)
			return
		}
		task.Date = newDate
		if _, err = dbt.UpdateTask(task); err != nil {
			retError(res, fmt.Sprintf("Tkd POST: Ошибка UpdateTask(): %s\n", err.Error()), http.StatusOK)
			return
		}
	}
	arrBytes, err := json.Marshal(bE)
	if err != nil {
		retError(res, fmt.Sprintf("Tkd POST: Ошибка json.Marshal(): %v\n", err), http.StatusOK)
		return
	}
	fmt.Printf("Tkd POST: ret json *%s*\n", string(arrBytes))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(arrBytes)
}

func TaskDoneHandle(res http.ResponseWriter, req *http.Request) {
	// одна задача (api/task/done)
	s := fmt.Sprintf("Получен запрос H: %s Path: %s M: %s", req.Host, req.URL.Path, req.Method)
	fmt.Println(s)
	if req.Method == "POST" {
		TaskDonePOSTHandle(res, req)
		return
	}
	retError(res, "Нужен только POST запрос", http.StatusOK)
}

func TasksHandle(res http.ResponseWriter, req *http.Request) {
	// много задач (api/tasks)
	s := fmt.Sprintf("Получен запрос H: %s Path: %s M: %s", req.Host, req.URL.Path, req.Method)
	fmt.Println(s)
	if req.Method == "GET" {
		TasksGETHandle(res, req)
		return
	}
	retError(res, "Ts Нужен только GET запрос", http.StatusOK)
}

func SignInPOSTHandle(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Запрос на авторизацию")
	var buf bytes.Buffer
	var pwds strcPwd
	// читаем тело запроса
	if _, err := buf.ReadFrom(req.Body); err != nil {
		retError(res, fmt.Sprintf("Si POST: Ошибка чтения тела запроса: %s\n", err.Error()), http.StatusOK)
		return
	}
	// десериализуем JSON в task
	if err := json.Unmarshal(buf.Bytes(), &pwds); err != nil {
		retError(res, fmt.Sprintf("Si POST: Ошибка десериализации: %s\n", err.Error()), http.StatusOK)
		return
	}
	// лог контроль
	fmt.Printf("Si POST: Unmarshal password *%s* \n", pwds.Password)
	// Функция должна сверять указанный пароль с хранимым в переменной окружения TODO_PASSWORD.
	// Если они совпадают, нужно сформировать JWT-токен и возвратить его в поле token JSON-объекта.
	envPassword := os.Getenv("TODO_PASSWORD")
	fmt.Printf("Si POST env password *%s*\n", envPassword)
	if envPassword == "" {
		// пароля в окружении нет. приваиваем свой
		envPassword = CPassword
	}
	if pwds.Password == envPassword {
		// при совпадении паролей SignInPOSHHandler возвращает в res token,
		// который frondend пишет в куки и который потом из куки используется для авторизации.
		// settings.Token нужно указывать только для тестирования алгоритма авторизации
		// Процесс смены пароля:
		// 1. Меняем CPassword
		// 2. Заходим из браузера с паролем из CPassword
		// 3. Из ответа сервера хэш копипастим в settings.Token для тетирования авторизации
		var tkn strcToken
		tkn.Token = JwtFromPass(envPassword)
		aBytes, _ := json.Marshal(tkn)
		// *** лог контроль
		fmt.Printf("Si POST: Marshal Token *%s*\n", string(aBytes))
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(aBytes)
		return
	}
	// Если пароль неверный или произошла ошибка, возвращается JSON c текстом ошибки в поле error.
	retError(res, "Пароль не верный", http.StatusUnauthorized)
}

func JwtFromPass(pass string) string {
	result := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(result[:])
}

func SignInHandle(res http.ResponseWriter, req *http.Request) {
	// авторизация (api/signin)
	s := fmt.Sprintf("Получен запрос H: %s Path: %s M: %s", req.Host, req.URL.Path, req.Method)
	fmt.Println(s)
	if req.Method == "POST" {
		SignInPOSTHandle(res, req)
		return
	}
	retError(res, "Нужен только POST запрос", http.StatusOK)
}
