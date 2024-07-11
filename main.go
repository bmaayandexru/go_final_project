package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bmaayandexru/go_final_project/db"
	"github.com/bmaayandexru/go_final_project/tests"
)

var mux *http.ServeMux

func main() {
	db.InitDBase()
	// лог-контроль
	fmt.Println("Запускаем сервер")
	mux = http.NewServeMux()
	// вешаем обработчик
	mux.HandleFunc("/m", mainHandle)
	//	webDir := "web/"
	//	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	// запуск файлового сервера в подкаталоге web
	mux.Handle("/", http.FileServer(http.Dir("web/")))
	// переменая tests.Port из settings.go
	// !!! заменить на переменную окружения TODO_PORT с помощью os.Getenv(key string) string
	// порт из settings.go
	settingsStrPort := fmt.Sprintf(":%d", tests.Port)
	// порт из переменной окружения. задание со *
	envStrPort := os.Getenv("TODO_PORT")
	// лог-контроль
	fmt.Printf("envPort *%s* settingdStrPort *%s* \n", envStrPort, settingsStrPort)
	fmt.Println("Set port from enviroment...")
	err := http.ListenAndServe(envStrPort, mux)

	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
}

func mainHandle(res http.ResponseWriter, req *http.Request) {
	// лог-контроль
	fmt.Println("Получен запрос")
	// запрос в строку
	s := fmt.Sprintf("Host: %s\nPath: %s\nMethod: %s", req.Host, req.URL.Path, req.Method)
	// лог-контроль
	fmt.Println(s)
	// отправка клиенту
	res.Write([]byte(s))
	// webDir := "web/" + req.URL.Path
	// mux.Handle("/", http.FileServer(http.Dir(webDir)))
}
