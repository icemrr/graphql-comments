package storage

import (
	"graphql-comments/internal/models"
	"testing"
)

func TestMemoryStorage_CreateAndGetPost(t *testing.T) {
	store := NewMemoryStorage()

	// 1. Создаем пост
	post := &models.Post{
		ID:      "post_1",
		Title:   "Мой пост",
		Content: "Текст поста",
	}

	err := store.CreatePost(post)
	if err != nil {
		t.Errorf("Ошибка при создании поста: %v", err)
	}

	// 2. Получаем пост
	saved, err := store.GetPost("post_1")
	if err != nil {
		t.Errorf("Ошибка при получении поста: %v", err)
	}

	// 3. Проверяем что тот же самый
	if saved.Title != "Мой пост" {
		t.Errorf("Ожидали 'Мой пост', получили '%s'", saved.Title)
	}
}

func TestMemoryStorage_CreateAndGetComment(t *testing.T) {
	store := NewMemoryStorage()

	// 1. Сначала создаем пост (обязательно)
	store.CreatePost(&models.Post{ID: "post_1", Title: "Пост", Content: "Контент"})

	// 2. Создаем комментарий
	comment := &models.Comment{
		ID:      "comment_1",
		PostID:  "post_1",
		Content: "Мой комментарий",
	}

	err := store.CreateComment(comment)
	if err != nil {
		t.Errorf("Ошибка при создании комментария: %v", err)
	}

	// 3. Получаем комментарии поста
	comments, err := store.GetCommentsByPostID("post_1")
	if err != nil {
		t.Errorf("Ошибка при получении комментариев: %v", err)
	}

	// 4. Проверяем что комментарий есть
	if len(comments) != 1 {
		t.Errorf("Ожидали 1 комментарий, получили %d", len(comments))
	}

	if comments[0].Content != "Мой комментарий" {
		t.Errorf("Ожидали 'Мой комментарий', получили '%s'", comments[0].Content)
	}
}
