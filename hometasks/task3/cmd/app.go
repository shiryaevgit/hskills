package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	serv "hskills_/hometasks/task3"
	"hskills_/hometasks/task3/config"
	"hskills_/hometasks/task3/pkg/handlers"
	"hskills_/hometasks/task3/pkg/repository"
	"hskills_/hometasks/task3/pkg/service"
)

/*
Собрать приложение через комманду: (запускаю из папки task3)
go build -o app cmd/app.g

Запустить приложение через консоль:
./app
*/

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("config.LoadEnv: %v", err)
	}

	port := env.HTTPPort
	srv := new(serv.Server)
	mux := http.NewServeMux()

	repo, err := repository.NewRepository(env.DatabaseFilepath)
	if err != nil {
		log.Fatalf("newRepository: %v", err)
	}
	metric := service.NewMetrics()
	handler := handlers.NewHandler(repo, metric)

	mux.HandleFunc("/healthcheck", handler.HandleGetHealthcheck)
	mux.HandleFunc("/redirect", handler.HandleGetRedirect)
	mux.HandleFunc("/values/{id}", handler.HandlePost)
	mux.HandleFunc("/values/", handler.HandleGet)

	//mux.HandleFunc("/values/", func(w http.ResponseWriter, r *http.Request) {
	//	if r.Method == http.MethodPost {
	//		handler.HandlePost(w, r)
	//	} else if r.Method == http.MethodGet {
	//		handler.HandleGet(w, r)
	//	}
	//})

	// Здесь зарегистрировал контекст который завершится в случае получения сигнала от операционной системы
	// Поэкспериментируй с ним и попробуй отправить разные сигналы процессу
	// Отправить SIGTERM процессу, id процесса можно узнать в консоли через команду ps:
	// kill -15 [pid]
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt) // Ctrl + C -> SIGINT; Есть еще SIGTERM,SIGKILL

	// Обработать ошибку корректного завершения сервера
	// errors.Is(err, http.ErrServerClosed)
	err = srv.Run(ctx, port, mux)
	if err != nil /* && */ {
		fmt.Println("Error run server:", err)
	} else {
		log.Printf("Server is running on http://127.0.0.1%s\n", port)
	}
}
