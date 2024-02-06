package handlers

import (
	"encoding/json"
	"fmt"
	"hskills_/pkg/models"
	"hskills_/pkg/repository"
	"hskills_/pkg/service"
	"log"
	"net/http"
	"net/url"
	"strconv"
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

// где его необходимо реализовывать/ использовать?
type UserRepository interface {
	CreatePost(newPost models.Post) error
	GetPost(id int) (*models.Post, error)
}

type Handler struct {
	repo    *repository.Repository
	metrics *service.Metrics
	// UserRepository ?
}

func NewHandler(repo *repository.Repository, metrics *service.Metrics) *Handler {
	return &Handler{repo: repo, metrics: metrics}
}

func (h *Handler) HandleGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodGet {
		load, err := h.metrics.GetCPULoad()
		if err != nil {
			log.Printf("HandleGetHealthcheck: GetCPULoad(): %v", err)
			http.Error(w, "internal error:", http.StatusInternalServerError)
		}
		metrics := models.HostMetric{
			CPULoad:      load,
			ThreadsCount: h.metrics.GetThreadCount(),
		}

		jsonData, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("HandleGetHealthcheck: json.Marshal(): %v", err)
			http.Error(w, "internal error:", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			log.Printf("HandleGetHealthcheck: w.Write(): %v", err)
			http.Error(w, "internal error:", http.StatusInternalServerError)
		}
	}
	fmt.Println("used handleGetHealthcheck")
}
func (h *Handler) HandleGetRedirect(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		targetURL := r.URL.Query().Get("url") //

		validTargetURL, err := url.Parse(targetURL)
		if err != nil {
			log.Printf("HandleGetRedirect: url.Parse(targetURL) %v", err)
			http.Error(w, "invalid link", http.StatusBadRequest)
		}
		http.Redirect(w, r, validTargetURL.String(), http.StatusFound)
	}
	fmt.Println("used HandleGetRedirect")
}
func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/values/"):]
	newPost := new(models.Post)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("HandlePost: strconv.Atoi(id): %v", err)
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	newPost.ID = idInt

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(newPost)
	if err != nil {
		log.Printf("HandlePost: decoder.Decode(newPost): %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	err = h.repo.CreatePost(*newPost)
	if err != nil {
		log.Printf("HandlePost: CreatePost(*newPost): %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/values/"):]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("HandleGet: Atoi(): %v", err)
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	post, err := h.repo.GetPost(idInt)
	if err != nil {
		log.Printf("HandleGet: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(post.Elements)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("HandleGet: Write(jsonData): %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
