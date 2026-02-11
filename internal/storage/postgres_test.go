package storage

import (
	"graphql-comments/internal/models"
	"testing"
)

func TestPostgresStorage_CreateAndGetPost(t *testing.T) {
	// 1. Подключаемся к тестовой БД
	dsn := "postgres://postgres:qwerty@localhost/comments_db_test?sslmode=disable"
	store, err := NewPostgresStorage(dsn)
	if err != nil {
		t.Skipf("PostgreSQL не доступен, пропускаем тест: %v", err)
	}
	defer store.Close()

	// 2. Очищаем таблицы перед тестом
	store.db.Exec("TRUNCATE posts, comments CASCADE;")

	// 3. Создаем пост
	post := &models.Post{
		ID:      "post_1",
		Title:   "Тестовый пост",
		Content: "Тестовое содержание",
	}

	err = store.CreatePost(post)
	if err != nil {
		t.Errorf("CreatePost вернул ошибку: %v", err)
	}

	// 4. Получаем пост
	saved, err := store.GetPost("post_1")
	if err != nil {
		t.Errorf("GetPost вернул ошибку: %v", err)
	}

	// 5. Проверяем
	if saved.Title != "Тестовый пост" {
		t.Errorf("Ожидали 'Тестовый пост', получили '%s'", saved.Title)
	}
}

func TestPostgresStorage_CreateAndGetComment(t *testing.T) {
	// 1. Подключаемся к тестовой БД
	dsn := "postgres://postgres:qwerty@localhost/comments_db_test?sslmode=disable"
	store, err := NewPostgresStorage(dsn)
	if err != nil {
		t.Skipf("PostgreSQL не доступен, пропускаем тест: %v", err)
	}
	defer store.Close()

	// 2. Очищаем таблицы
	store.db.Exec("TRUNCATE posts, comments CASCADE;")

	// 3. Сначала создаем пост
	post := &models.Post{ID: "post_1", Title: "Пост", Content: "Контент"}
	store.CreatePost(post)

	// 4. Создаем комментарий
	comment := &models.Comment{
		ID:      "comment_1",
		PostID:  "post_1",
		Content: "Тестовый комментарий",
	}

	err = store.CreateComment(comment)
	if err != nil {
		t.Errorf("CreateComment вернул ошибку: %v", err)
	}

	// 5. Получаем комментарии
	comments, err := store.GetCommentsByPostID("post_1")
	if err != nil {
		t.Errorf("GetCommentsByPostID вернул ошибку: %v", err)
	}

	// 6. Проверяем
	if len(comments) != 1 {
		t.Errorf("Ожидали 1 комментарий, получили %d", len(comments))
	}

	if comments[0].Content != "Тестовый комментарий" {
		t.Errorf("Ожидали 'Тестовый комментарий', получили '%s'", comments[0].Content)
	}
}
