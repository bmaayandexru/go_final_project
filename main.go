package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bmaayandexru/go_final_project/dbt"
	"github.com/bmaayandexru/go_final_project/tests"
)

var mux *http.ServeMux

func main() {
	dbt.InitDBase()
	// лог-контроль
	fmt.Println("Запускаем сервер")
	mux = http.NewServeMux()
	// вешаем отладочный обработчик
	mux.HandleFunc("/m", mainHandle)
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

// отладочный обработчик. убрать
func mainHandle(res http.ResponseWriter, req *http.Request) {
	// лог-контроль
	fmt.Println("Получен запрос")
	// запрос в строку
	s := fmt.Sprintf("Host: %s\nPath: %s\nMethod: %s", req.Host, req.URL.Path, req.Method)
	// лог-контроль
	fmt.Println(s)
	// отправка клиенту
	res.Write([]byte(s))
}
