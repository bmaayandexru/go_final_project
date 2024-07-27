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

/*
	type nextDate struct {
		date   string
		repeat string
		want   string
	}
*/
func main() {
	// now=20240126
	//{"20240202", "d 30", `20240303`},//*
	//{"20240228", "d 1", "20240229"},
	//{"20240126", "m 25,26,7", "20240207"}
	//{"20230126", "w 4,5", "20240201"}
	/*
		dnow, e := time.Parse("20060102", "20240126")
		nd := nextDate{"20230126", "w 4,5", "20240201"}
		s, e := handlers.NextDate(dnow, nd.date, nd.repeat)
		fmt.Printf("retstr *%s* want *%s* err *%v*\n", s, nd.want, e)
		return
	*/
	//***
	dbt.InitDBase()
	mux = http.NewServeMux()

	mux.HandleFunc("/api/nextdate", handlers.NextDateHandle)
	mux.HandleFunc("/api/task", handlers.TaskHandle)
	mux.HandleFunc("/api/task/done", handlers.TaskDoneHandle)
	mux.HandleFunc("/api/tasks", handlers.TasksHandle)
	mux.Handle("/", http.FileServer(http.Dir("web/")))
	strPort := defStrPort()
	err := http.ListenAndServe(strPort, mux)
	if err != nil {
		panic(err)
	}
}

func defStrPort() string {
	defPort := "7540"
	// переменая tests.Port из settings.go
	settingsStrPort := fmt.Sprintf("%d", tests.Port)
	// лог-контроль
	if settingsStrPort != "" {
		defPort = settingsStrPort
	}
	envStrPort := os.Getenv("TODO_PORT")
	if envStrPort != "" {
		defPort = envStrPort
	}
	fmt.Printf("Set port %s \n", defPort)
	return ":" + defPort
}
