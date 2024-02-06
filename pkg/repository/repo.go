package repository

import (
	"encoding/json"
	"fmt"
	"hskills_/pkg/models"
	"io"
	"os"
	"sync"
)

/*
	Место реализации
*/

type Repository struct {
	file *os.File
	mu   sync.Mutex
	path string
}

func NewRepository(path string) (*Repository, error) {
	file, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("NewRepository: %w", err)
	}

	return &Repository{file: file, path: path}, nil
}

//func NewRepository(path string) (*Repository, error) {
//	openedFile, err := os.Open(path)
//	if err == nil {
//		return &Repository{file: openedFile, path: path}, nil
//	}
//
//	createFile, err := os.Create(path)
//	if err != nil {
//		return nil, fmt.Errorf("NewRepository: %w", err)
//	}
//	return &Repository{file: createFile, path: path}, nil
//}

func (r *Repository) CreatePost(newPost models.Post) error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	posts, err := r.GetAllPosts()
	if err != nil {
		return fmt.Errorf("CreatePost: %v", err)
	}

	for _, post := range posts {
		if post.ID == newPost.ID {
			return fmt.Errorf("post with ID:%d already exist", newPost.ID)
		}
	}

	posts = append(posts, newPost)
	jsonData, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("CreatePost: %v", err)
	}

	err = os.WriteFile(r.path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("CreatePost: %v", err)
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
	return nil, fmt.Errorf("GetPost: there is no post with id: %v", id)
}

func (r *Repository) GetAllPosts() ([]models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var posts []models.Post

	_, err := r.file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("GetAllPosts: %v", err)
	}

	data, err := io.ReadAll(r.file)
	if len(data) == 0 {
		fmt.Println("post list is empty")
		return posts, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetAllPosts: %w", err)
	}

	err = json.Unmarshal(data, &posts)
	if err != nil {
		return nil, fmt.Errorf("GetAllPosts: %w", err)
	}

	return posts, nil
}
