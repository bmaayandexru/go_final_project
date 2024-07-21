package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bmaayandexru/go_final_project/dbt"
	"github.com/bmaayandexru/go_final_project/handlers"
	"github.com/bmaayandexru/go_final_project/tests"
)

var mux *http.ServeMux

func main() {
	/*
		//d, e := time.Parse("20060102", "20240126")
		d := time.Now()
		s, e := handlers.NextDate(d, "20240409", "m 1,7 2,12")
		fmt.Printf("retstr *%s* err *%v*\n", s, e)
		return
	*/
	//***
	dbt.InitDBase()
	// лог-контроль
	fmt.Println("Запускаем сервер")
	mux = http.NewServeMux()
	// вешаем отладочный обработчик
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandle)
	mux.HandleFunc("/api/task", handlers.TaskHandle)
	// запуск файлового сервера в подкаталоге web
	mux.Handle("/", http.FileServer(http.Dir("web/")))
	// определение порта прослушки
	strPort := defStrPort()
	err := http.ListenAndServe(strPort, mux)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
}

func defStrPort() string {
	defPort := "7540"
	// переменая tests.Port из settings.go
	settingsStrPort := fmt.Sprintf("%d", tests.Port)
	// лог-контроль
	fmt.Printf("settingdStrPort *%s* \n", settingsStrPort)
	if settingsStrPort != "" {
		// значение не пустое. переобределяем
		defPort = settingsStrPort
	}
	// порт из переменной окружения TODO_PORT. задание со *
	envStrPort := os.Getenv("TODO_PORT")
	// лог-контроль
	fmt.Printf("envPort *%s* \n", envStrPort)
	if envStrPort != "" {
		// значение не пустое. переобределяем
		defPort = envStrPort
	}
	// итоговое значение
	fmt.Printf("Set port %s \n", defPort)
	return ":" + defPort
}
