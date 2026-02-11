package storage

import "graphql-comments/internal/models"

// Storage - это интерфейс, который определяет все методы,
// которые должно поддерживать наше хранилище (in-memory или postgres)
type Storage interface {
	// Методы для работы с постами
	CreatePost(post *models.Post) error
	GetPost(id string) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	DeletePost(id string) error

	// Методы для работы с комментариями
	CreateComment(comment *models.Comment) error
	GetComment(id string) (*models.Comment, error)
	GetCommentsByPostID(postID string) ([]*models.Comment, error)
	DeleteComment(id string) error
}
