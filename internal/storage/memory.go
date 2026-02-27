package storage

import (
	"errors"
	"graphql-comments/internal/models"
	"sync"
)

// MemoryStorage - реализация Storage, которая хранит данные в памяти
type MemoryStorage struct {
	mu       sync.RWMutex               
	posts    map[string]*models.Post    
	comments map[string]*models.Comment 
}

// NewMemoryStorage создает новый экземпляр MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts:    make(map[string]*models.Post),
		comments: make(map[string]*models.Comment),
	}
}

// CreatePost создает новый пост
func (s *MemoryStorage) CreatePost(post *models.Post) error {
	s.mu.Lock()        
	defer s.mu.Unlock() 

	// Проверяем, существует ли уже пост
	if _, exists := s.posts[post.ID]; exists {
		return errors.New("пост уже существует")
	}

	// Инициализируем Comments слайс
	if post.Comments == nil {
		post.Comments = []*models.Comment{}
	}

	// Сохраняем пост в мапе
	s.posts[post.ID] = post
	return nil
}

// GetPost возвращает пост по ID
func (s *MemoryStorage) GetPost(id string) (*models.Post, error) {
	s.mu.RLock()        
	defer s.mu.RUnlock() 

	post, exists := s.posts[id]
	if !exists {
		return nil, errors.New("пост не найден")
	}

	// Создаем копию поста
	// и инициализируем Comments если нужно
	postCopy := *post
	if postCopy.Comments == nil {
		postCopy.Comments = []*models.Comment{}
	}

	return &postCopy, nil
}

// GetAllPosts возвращает все посты
func (s *MemoryStorage) GetAllPosts() ([]*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Создаем слайс для результата
	posts := make([]*models.Post, 0, len(s.posts))

	// Копируем все посты
	for _, post := range s.posts {
		postCopy := *post
		if postCopy.Comments == nil {
			postCopy.Comments = []*models.Comment{}
		}
		posts = append(posts, &postCopy)
	}

	return posts, nil
}

// DeletePost удаляет пост по ID
func (s *MemoryStorage) DeletePost(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.posts[id]; !exists {
		return errors.New("пост не найден")
	}

	delete(s.posts, id)
	return nil
}

// CreateComment создает новый комментарий
func (s *MemoryStorage) CreateComment(comment *models.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем существование комментария
	if _, exists := s.comments[comment.ID]; exists {
		return errors.New("комментарий уже существует")
	}

	// Проверяем, существует ли пост, к которому добавляем комментарий
	if _, exists := s.posts[comment.PostID]; !exists {
		return errors.New("пост не найден")
	}

	// Если есть ParentID, проверяем существование родительского комментария
	if comment.ParentID != nil {
		if _, exists := s.comments[*comment.ParentID]; !exists {
			return errors.New("родительский комментарий не найден")
		}
	}

	// Инициализируем Replies слайс
	if comment.Replies == nil {
		comment.Replies = []*models.Comment{}
	}

	// Сохраняем комментарий
	s.comments[comment.ID] = comment
	return nil
}

// GetComment возвращает комментарий по ID
func (s *MemoryStorage) GetComment(id string) (*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comment, exists := s.comments[id]
	if !exists {
		return nil, errors.New("комментарий не найден")
	}

	// Создаем копию
	commentCopy := *comment
	if commentCopy.Replies == nil {
		commentCopy.Replies = []*models.Comment{}
	}

	return &commentCopy, nil
}

// GetCommentsByPostID возвращает ВСЕ комментарии для указанного поста
func (s *MemoryStorage) GetCommentsByPostID(postID string) ([]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Проверяем существование поста
	if _, exists := s.posts[postID]; !exists {
		return nil, errors.New("пост не найден")
	}

	// Собираем все комментарии для этого поста
	var comments []*models.Comment
	for _, comment := range s.comments { // бежим по все мапе
		if comment.PostID == postID {
			// Создаем копию комментария
			commentCopy := *comment
			if commentCopy.Replies == nil {
				commentCopy.Replies = []*models.Comment{}
			}
			comments = append(comments, &commentCopy)
		}
	}
	// вОзвращаем плоским списком
	return comments, nil
}

// DeleteComment удаляет комментарий по ID и все его ответы рекурсивно
func (s *MemoryStorage) DeleteComment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем существование комментария
	if _, exists := s.comments[id]; !exists {
		return errors.New("комментарий не найден")
	}

	// Рекурсивно удаляем все дочерние комментарии
	s.deleteCommentRecursive(id)

	return nil
}

// deleteCommentRecursive удаляет комментарий и все его ответы
func (s *MemoryStorage) deleteCommentRecursive(id string) {
	// Удаляем текущий комментарий
	delete(s.comments, id)

	// Ищем и удаляем все комментарии, у которых этот комментарий - родитель
	for commentID, comment := range s.comments {
		if comment.ParentID != nil && *comment.ParentID == id {
			s.deleteCommentRecursive(commentID)
		}
	}
}

var _ Storage = (*MemoryStorage)(nil)
