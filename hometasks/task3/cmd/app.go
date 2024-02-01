package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	serv "hskills_/hometasks/task3"
	"hskills_/hometasks/task3/pkg/handlers"
	"hskills_/hometasks/task3/pkg/repository"
	"hskills_/hometasks/task3/pkg/service"
)

/*
type Environment struct {
	HTTPPort string
	DatabaseFilepath string
}
*/

/*
Собрать приложение через комманду: (запускаю из папки task3)
go build -o app cmd/app.g

Запустить приложение через консоль:
./app
*/

func main() {
	port := "8080"
	srv := new(serv.Server)
	mux := http.NewServeMux() // Создаем новый мультиплексор - обработчик запросов/путей

	repo, err := repository.NewRepository("./db.json")
	if err != nil {
		/*
			Не принято (от слова совсем) читать значение из функции в случае если функция вернула ошибку
		*/
		log.Fatalln("ошибка при создании репозитория", err)
	}
	metric := service.NewMetrics()
	handler := handlers.NewHandler(repo, metric)

	// Создаем обработчики для конкретных путей:
	mux.HandleFunc("/healthcheck", handler.HandleGetHealthcheck) // передаем в мультиплексор путь и функции для обработки
	mux.HandleFunc("/redirect", handler.HandleGetRedirect)
	mux.HandleFunc("/values/{id}", handler.HandlePost)
	mux.HandleFunc("/values/", handler.HandleGet)

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
