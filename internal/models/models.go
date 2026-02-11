package models

// Post представляет собой пост/статью, к которой пишутся комментарии
type Post struct {
	ID       string     `json:"id"`       // Уникальный идентификатор поста
	Title    string     `json:"title"`    // Заголовок поста
	Content  string     `json:"content"`  // Содержимое поста
	Comments []*Comment `json:"comments"` // Список комментариев первого уровня
}

// Comment представляет собой комментарий, который может быть ответом на пост или другой комментарий
type Comment struct {
	ID       string     `json:"id"`       // Уникальный идентификатор комментария
	PostID   string     `json:"postId"`   // ID поста, к которому относится комментарий
	ParentID *string    `json:"parentId"` // ID родительского комментария (nil для комментариев первого уровня)
	Content  string     `json:"content"`  // Текст комментария
	Replies  []*Comment `json:"replies"`  // Дочерние комментарии (ответы на этот комментарий)
}

// CreatePostInput - входные данные для создания поста
type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateCommentInput - входные данные для создания комментария
type CreateCommentInput struct {
	PostID   string  `json:"postId"`   // Обязательно
	ParentID *string `json:"parentId"` // Необязательно, может быть nil
	Content  string  `json:"content"`  // Обязательно
}
