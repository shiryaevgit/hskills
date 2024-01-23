package repository

import (
	"encoding/json"
	"fmt"
	"hskills_/hometasks/task3/pkg/models"
	"os"
)

type Repository struct {
	file *os.File
}

func NewRepository(name string) (*Repository, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("error create database: %w", err)
	}
	return &Repository{file: file}, nil
}

func (r *Repository) CreatePost(newPost models.Post) error {
	posts, err := r.GetAllPosts()
	if err != nil {
		return err
	}

	for _, post := range posts {
		if post.ID == newPost.ID {
			return fmt.Errorf("post with ID%d already exist", newPost.ID)
		}
	}

	posts = append(posts, newPost)
	jsonData, err := json.Marshal(posts)

	err = os.WriteFile("db.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func (r *Repository) GetPost(id int) (*models.Post, error) {
	posts, err := r.GetAllPosts()
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		if post.ID == id {
			return &post, nil
		}
	}
	return nil, fmt.Errorf("there is no post with this id:%v", id)
}

func (r *Repository) GetAllPosts() ([]models.Post, error) {

	data, err := os.ReadFile("db.json")
	if err != nil {
		return nil, fmt.Errorf("error ReadFile:%w", err)
	}

	var posts []models.Post
	err = json.Unmarshal(data, &posts)
	if err != nil {
		return nil, fmt.Errorf("error Unmarshal:%w", err)
	}
	fmt.Println(posts)
	return posts, err
}
