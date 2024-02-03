package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"hskills_/hometasks/task3/pkg/models"
	"hskills_/hometasks/task3/pkg/repository"
	"hskills_/hometasks/task3/pkg/service"
)

/*
	Место применения (Handler)
*/

/*  Пример интерфейса репозитория.
Заполняем только методы которые нужны в сервисе

type Repository interface {
	CreatePost(newPost models.Post) error
}
*/

type UserRepository interface {
	CreatePost(newPost models.Post) error
	GetPost(id int) (*models.Post, error)
}

type Handler struct {
	repo    *repository.Repository
	metrics *service.Metrics
}

func NewHandler(repo *repository.Repository, metrics *service.Metrics) *Handler {
	return &Handler{repo: repo, metrics: metrics}
}
func (h *Handler) HandleGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		load, err := h.metrics.GetCPULoad()
		if err != nil {
			// Ошибка не обработана (сделать log + http.Error)
			/*
				Ответ пользователю: { "code": 500, "status": "internal error" }
				Лог: детальный вывод информации. Ошибки прокидываем наверх, логируем
			*/
			return
		}
		metrics := models.HostMetric{
			CPULoad:      load,
			ThreadsCount: h.metrics.GetThreadCount(),
		}

		jsonData, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, "ошибка формирования json", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, "ошибка формирования json", http.StatusInternalServerError)
			return
		}
	}
	fmt.Println("used handleGetHealthcheck")
}
func (h *Handler) HandleGetRedirect(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		targetURL := r.URL.Query().Get("url") //

		validTargetURL, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, "invalid link", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, validTargetURL.String(), http.StatusFound)
	}

}
func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("used Post")

	if r.Method == http.MethodPost {
		id := r.URL.Path[len("/values/"):]
		newPost := new(models.Post)

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}
		newPost.ID = idInt

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(newPost)
		if err != nil {
			http.Error(w, "error decoding json", http.StatusBadRequest)
			return
		}

		err = h.repo.CreatePost(*newPost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Path[len("/values/"):]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}

		post, err := h.repo.GetPost(idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonData, err := json.Marshal(post.Elements)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, "formation error json", http.StatusInternalServerError)
		}

		fmt.Println("used Get")
	}
}
