package gql

import (
	"strconv"
	"sync"

	"graphql-comments/internal/models"
	"graphql-comments/internal/storage"

	"github.com/graphql-go/graphql"
)

// ResolverContext хранит зависимости для резолверов
type ResolverContext struct {
	Storage        storage.Storage
	mu             sync.Mutex
	postCounter    int
	commentCounter int
}

// PostsResolver возвращает все посты
func (r *ResolverContext) PostsResolver(p graphql.ResolveParams) (interface{}, error) {
	posts, err := r.Storage.GetAllPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// CreatePostResolver создает новый пост
func (r *ResolverContext) CreatePostResolver(p graphql.ResolveParams) (interface{}, error) {
	title, _ := p.Args["title"].(string)
	content, _ := p.Args["content"].(string)

	post := &models.Post{
		ID:       r.generatePostID(),
		Title:    title,
		Content:  content,
		Comments: []*models.Comment{},
	}

	err := r.Storage.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// DeletePostResolver удаляет пост по ID
func (r *ResolverContext) DeletePostResolver(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(string)

	err := r.Storage.DeletePost(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateCommentResolver создает новый комментарий
func (r *ResolverContext) CreateCommentResolver(p graphql.ResolveParams) (interface{}, error) {
	// Получаем input объект
	input, _ := p.Args["input"].(map[string]interface{})

	postID, _ := input["postId"].(string)
	content, _ := input["content"].(string)

	// ParentID может быть nil или строкой
	var parentID *string
	if parentArg, ok := input["parentId"].(string); ok {
		parentID = &parentArg
	}

	comment := &models.Comment{
		ID:       r.generateCommentID(),
		PostID:   postID,
		ParentID: parentID,
		Content:  content,
		Replies:  []*models.Comment{},
	}

	err := r.Storage.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// DeleteCommentResolver удаляет комментарий по ID
func (r *ResolverContext) DeleteCommentResolver(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(string)

	err := r.Storage.DeleteComment(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CommentsResolver возвращает комментарии для поста (плоский список)
func (r *ResolverContext) CommentsResolver(p graphql.ResolveParams) (interface{}, error) {
	// p.Source содержит родительский объект (Post)
	post, ok := p.Source.(*models.Post)
	if !ok {
		return nil, nil
	}

	// Получаем все комментарии для этого поста
	comments, err := r.Storage.GetCommentsByPostID(post.ID)
	if err != nil {
		return nil, err
	}

	// Преобразуем плоский список в дерево
	return r.buildCommentTree(comments), nil
}

// RepliesResolver возвращает ответы на комментарий
func (r *ResolverContext) RepliesResolver(p graphql.ResolveParams) (interface{}, error) {
	// p.Source содержит родительский объект (Comment)
	comment, ok := p.Source.(*models.Comment)
	if !ok {
		return nil, nil
	}

	// Для рекурсивного построения дерева
	// В реальности нужно получать все комментарии и фильтровать
	// Но мы сделаем проще - вернем уже готовое поле Replies
	return comment.Replies, nil
}

// buildCommentTree преобразует плоский список комментариев в дерево
func (r *ResolverContext) buildCommentTree(comments []*models.Comment) []*models.Comment {
	// Создаем мапу для быстрого доступа: ID комментария -> комментарий
	commentMap := make(map[string]*models.Comment)

	// Копируем комментарии в мапу
	for _, comment := range comments {
		commentCopy := *comment
		commentCopy.Replies = []*models.Comment{} // Инициализируем пустой слайс
		commentMap[commentCopy.ID] = &commentCopy
	}

	// Собираем дерево
	var rootComments []*models.Comment
	for _, comment := range commentMap {
		if comment.ParentID == nil {
			// Корневой комментарий (не имеет родителя)
			rootComments = append(rootComments, comment)
		} else {
			// Ответ на другой комментарий
			parent, exists := commentMap[*comment.ParentID]
			if exists {
				parent.Replies = append(parent.Replies, comment)
			}
		}
	}

	return rootComments
}

// generatePostID генерирует уникальный ID для поста
func (r *ResolverContext) generatePostID() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.postCounter++
	return "post_" + strconv.Itoa(r.postCounter)
}

// generateCommentID генерирует уникальный ID для комментария
func (r *ResolverContext) generateCommentID() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.commentCounter++
	return "comment_" + strconv.Itoa(r.commentCounter)
}
