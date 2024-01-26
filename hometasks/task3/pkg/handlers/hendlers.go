package handlers

import (
	"encoding/json"
	"fmt"
	"hskills_/hometasks/task3/pkg/models"
	"hskills_/hometasks/task3/pkg/repository"
	"hskills_/hometasks/task3/pkg/service"
	"net/http"
	"net/url"
	"strconv"
)

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
	fmt.Println("used HandleGetRedirect")
}
func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.URL.Path[len("/values/"):]

		newPost := new(models.Post)

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid ID", 400)
			return
		}
		newPost.ID = idInt

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(newPost)
		if err != nil {
			http.Error(w, "error decoding json", 400)
			return
		}

		err = h.repo.CreatePost(*newPost)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	fmt.Println("used handlePost")
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Path[len("/values/"):]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid ID", 400)
			return
		}

		post, err := h.repo.GetPost(idInt)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		jsonData, err := json.Marshal(post.Elements)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, "formation error json", http.StatusInternalServerError)
			return
		}
	}

}
