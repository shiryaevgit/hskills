package handlers

import (
	"encoding/json"
	"fmt"
	"hskills_/hometasks/task3/pkg/models"
	"hskills_/hometasks/task3/pkg/service"
	"os"

	"net/http"
	"strconv"
)

func HandleGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		load, err := service.GetCPULoad()
		if err != nil {
			return
		}
		metrics := models.HostMetric{
			CPULoad:      load,
			ThreadsCount: service.GetThreadCount(),
		}

		jsonData, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, "ошибка формирования json", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, "ошибка формирования json", http.StatusInternalServerError)
		}
	}
	fmt.Println("used handleGetHealthcheck")
}
func HandleGetRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		targetURL := r.URL.Query().Get("URL") //

		if targetURL != "" {
			http.Redirect(w, r, targetURL, http.StatusFound)
		}
	}
	fmt.Println("used handleGetRedirect")
}
func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.URL.Path[len("/values/"):]

		newPost := new(models.Post)

		idString, err := strconv.Atoi(id)
		newPost.ID = idString

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(newPost)
		if err != nil {
			http.Error(w, "ошибка декодирования json", 400)
		}

		jsonData, err := json.Marshal(newPost)

		// Далее надо вызывать методы репозитория)?
		err = os.WriteFile("./db.json", jsonData, 0644)
		if err != nil {
			fmt.Errorf("error writing to file: %w", err)
		}

		fmt.Printf("ID: %s, Elements: %v\n", id, newPost.Elements)
	}

	fmt.Println("used handlePost")
}
func HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Path[len("/values/"):]

		//не понимаю как обратится к методам репозитория чтобы вызвать
		//GetAllPosts и проверить есть ли в наличии

	}
}
