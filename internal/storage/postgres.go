package storage

import (
	"database/sql"
	"fmt"
	"graphql-comments/internal/models"

	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

// PostgresStorage реализация Storage для PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage создает новое подключение к PostgreSQL
func NewPostgresStorage(dataSourceName string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// Close закрывает подключение к БД
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

// CreatePost создает новый пост в БД
func (s *PostgresStorage) CreatePost(post *models.Post) error {
	query := `INSERT INTO posts (id, title, content) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, post.ID, post.Title, post.Content)
	return err
}

// GetPost возвращает пост по ID из БД
func (s *PostgresStorage) GetPost(id string) (*models.Post, error) {
	query := `SELECT id, title, content FROM posts WHERE id = $1`
	row := s.db.QueryRow(query, id)

	post := &models.Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, err
	}

	post.Comments = []*models.Comment{}
	return post, nil
}

// GetAllPosts возвращает все посты из БД
func (s *PostgresStorage) GetAllPosts() ([]*models.Post, error) {
	query := `SELECT id, title, content FROM posts ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			return nil, err
		}
		post.Comments = []*models.Comment{}
		posts = append(posts, post)
	}

	return posts, nil
}

// DeletePost удаляет пост по ID из БД
func (s *PostgresStorage) DeletePost(id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

// CreateComment создает новый комментарий в БД
func (s *PostgresStorage) CreateComment(comment *models.Comment) error {
	var query string
	var args []interface{}

	if comment.ParentID != nil {
		query = `INSERT INTO comments (id, post_id, parent_id, content) VALUES ($1, $2, $3, $4)`
		args = []interface{}{comment.ID, comment.PostID, *comment.ParentID, comment.Content}
	} else {
		query = `INSERT INTO comments (id, post_id, content) VALUES ($1, $2, $3)`
		args = []interface{}{comment.ID, comment.PostID, comment.Content}
	}

	_, err := s.db.Exec(query, args...)
	return err
}

// GetComment возвращает комментарий по ID из БД
func (s *PostgresStorage) GetComment(id string) (*models.Comment, error) {
	query := `SELECT id, post_id, parent_id, content FROM comments WHERE id = $1`
	row := s.db.QueryRow(query, id)

	comment := &models.Comment{}
	var parentID sql.NullString
	err := row.Scan(&comment.ID, &comment.PostID, &parentID, &comment.Content)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("comment not found")
	}
	if err != nil {
		return nil, err
	}

	if parentID.Valid {
		parentIDStr := parentID.String
		comment.ParentID = &parentIDStr
	}

	comment.Replies = []*models.Comment{}
	return comment, nil
}

// GetCommentsByPostID возвращает все комментарии для поста из БД
func (s *PostgresStorage) GetCommentsByPostID(postID string) ([]*models.Comment, error) {
	query := `SELECT id, post_id, parent_id, content FROM comments WHERE post_id = $1 ORDER BY created_at`
	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		var parentID sql.NullString

		if err := rows.Scan(&comment.ID, &comment.PostID, &parentID, &comment.Content); err != nil {
			return nil, err
		}

		if parentID.Valid {
			parentIDStr := parentID.String
			comment.ParentID = &parentIDStr
		}

		comment.Replies = []*models.Comment{}
		comments = append(comments, comment)
	}

	return comments, nil
}

// DeleteComment удаляет комментарий по ID из БД
func (s *PostgresStorage) DeleteComment(id string) error {
	query := `DELETE FROM comments WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}
