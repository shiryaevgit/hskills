package main

import (
	"context"
	"errors"
	"fmt"
	serv "hskills_"
	"hskills_/config"
	"hskills_/pkg/handlers"
	"hskills_/pkg/repository"
	"hskills_/pkg/service"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	pid := os.Getpid()
	fmt.Println("Current PID:", pid)

	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("config.LoadEnv: %v", err)
	}

	fileLog, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("create error.log: %v", err)
	}
	log.SetOutput(fileLog)

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
	mux.HandleFunc("/values/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.HandlePost(w, r)
		} else if r.Method == http.MethodGet {
			handler.HandleGet(w, r)
		}
	})

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt) // Ctrl + C -> SIGINT; Есть еще SIGTERM,SIGKILL

	err = srv.Run(ctx, port, mux)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Run():", err)
		} else {
			fmt.Printf("Run(): %v", err)
		}
	} else {
		log.Printf("Server is running on http://127.0.0.1%s\n", port)
	}

}
