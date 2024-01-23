package main

import (
	"fmt"
	_ "github.com/lib/pq"
	serv "hskills_/hometasks/task3"
	"hskills_/hometasks/task3/pkg/handlers"
	"hskills_/hometasks/task3/pkg/repository"
	"log"
	"net/http"
)

func main() {
	port := "8080"
	srv := new(serv.Server)
	mux := http.NewServeMux() // Создаем новый мультиплексор - обработчик запросов/путей

	_, err := repository.NewRepository("./db.json")
	if err != nil {
		fmt.Println("ошибка при создании репозитория")
	}
	// Создаем обработчики для конкретных путей:
	mux.HandleFunc("/healthcheck", handlers.HandleGetHealthcheck) // передаем в мультиплексор путь и функции для обработки
	mux.HandleFunc("/values/", handlers.HandlePost)
	mux.HandleFunc("/redirect", handlers.HandleGetRedirect)
	//mux.HandleFunc("/values/", handlers.HandleGet)

	err = srv.Run(port, mux)
	if err != nil {
		fmt.Println("Error run server:", err)
	} else {
		log.Printf("Server is running on http://127.0.0.1%s\n", port)
	}

}
