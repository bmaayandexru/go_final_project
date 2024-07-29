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
	dbt.InitDBase()
	mux = http.NewServeMux()

	mux.HandleFunc("/api/nextdate", handlers.NextDateHandle)
	mux.HandleFunc("/api/task", auth(handlers.TaskHandle))
	mux.HandleFunc("/api/task/done", auth(handlers.TaskDoneHandle))
	mux.HandleFunc("/api/tasks", auth(handlers.TasksHandle))
	mux.HandleFunc("/api/signin", handlers.SignInHandle)
	mux.Handle("/", http.FileServer(http.Dir("web/")))
	strPort := defStrPort()
	err := http.ListenAndServe(strPort, mux)
	if err != nil {
		panic(err)
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		fmt.Printf("env password *%s*\n", pass)
		// авторизация будет проверяться только при наличии TODO_PASSWORD
		// иначе не будет
		if len(pass) == 0 {
			// это чтоб работало без переменной окружения
			pass = handlers.CPassword
		}
		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			var valid bool
			jwtp := handlers.JwtFromPass(pass)
			fmt.Printf("auth: jwt pass *%s*\n", jwtp)
			// здесь код для валидации и проверки JWT-токена
			fmt.Println("auth: cookie jwt", jwt)
			valid = (jwt == jwtp)
			//valid = (jwt == tests.Token)
			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
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
	fmt.Printf("env port *%s* \n", envStrPort)
	if envStrPort != "" {
		defPort = envStrPort
	}
	fmt.Printf("Set port %s \n", defPort)
	return ":" + defPort
}
