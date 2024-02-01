package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"hskills_/hometasks/task3/pkg/models"
)

/*
	Место реализации
*/

type Repository struct {
	file *os.File
	mu   sync.Mutex
}

func NewRepository(name string) (*Repository, error) {
	openedFile, err := os.Open(name)
	if err == nil {
		return &Repository{file: openedFile}, nil
	}

	createFile, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("NewRepository: create database: %w", err)
	}
	return &Repository{file: createFile}, nil
}

func (r *Repository) CreatePost(newPost models.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	posts, err := r.GetAllPosts()
	if err != nil {
		return fmt.Errorf("CreatePost: %w", err)
	}

	for _, post := range posts {
		if post.ID == newPost.ID {
			return fmt.Errorf("post with ID:%d already exist", newPost.ID)
		}
	}

	posts = append(posts, newPost)
	jsonData, err := json.Marshal(posts)

	err = os.WriteFile("db.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("CreatePost: os.WriteFile() %w", err)
	}
	return nil
}

func (r *Repository) GetPost(id int) (*models.Post, error) {
	posts, err := r.GetAllPosts()
	if err != nil {
		return nil, fmt.Errorf("GetPost: GetAllPosts() %v", err)
	}
	for _, post := range posts {
		if post.ID == id {
			return &post, nil
		}
	}
	return nil, fmt.Errorf("there is no post with this id:%v", id)
}

func (r *Repository) GetAllPosts() ([]models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var posts []models.Post

	_, err := r.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("GetAllPosts: file.Seek() %w", err)
	}

	data, err := io.ReadAll(r.file)
	if len(data) == 0 {
		fmt.Println("post list is empty")
		return posts, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetAllPosts: io.ReadAll() %w", err)
	}

	err = json.Unmarshal(data, &posts)
	if err != nil {
		return nil, fmt.Errorf("GetAllPosts: json.Unmarshal() %w", err)
	}
	return posts, nil
}
