package gql

import (
	"graphql-comments/internal/models"
	"graphql-comments/internal/storage"
	"testing"
)


func TestBuildCommentTree_Empty(t *testing.T) {
	// Создаем резолвер с in-memory хранилищем
	resolver := &ResolverContext{
		Storage: storage.NewMemoryStorage(),
	}

	// Пустой список комментариев - ситуация когда у поста нет комментариев
	var comments []*models.Comment

	// Вызываем тестируемую функцию
	result := resolver.buildCommentTree(comments)

	// Проверяем: результат должен быть пустым массивом
	if len(result) != 0 {
		t.Errorf("Ожидали пустой массив, получили %d", len(result))
	}
}

// TestBuildCommentTree_RootComments - проверяет что корневые комментарии
// (parentId = null) правильно определяются и возвращаются на верхнем уровне
func TestBuildCommentTree_RootComments(t *testing.T) {
	resolver := &ResolverContext{
		Storage: storage.NewMemoryStorage(),
	}

	// Создаем два комментария без parentId - они должны стать корневыми
	comments := []*models.Comment{
		{ID: "comment_1", PostID: "post_1", Content: "Коммент 1", Replies: []*models.Comment{}},
		{ID: "comment_2", PostID: "post_1", Content: "Коммент 2", Replies: []*models.Comment{}},
	}

	// buildCommentTree должна вернуть оба комментария на верхнем уровне
	result := resolver.buildCommentTree(comments)

	// Проверяем что оба комментария стали корневыми
	if len(result) != 2 {
		t.Errorf("Ожидали 2 корневых комментария, получили %d", len(result))
	}
}


func TestBuildCommentTree_WithReplies(t *testing.T) {
	resolver := &ResolverContext{
		Storage: storage.NewMemoryStorage(),
	}

	// Создаем указатели на ID для parentId
	parentID1 := "comment_1"
	parentID2 := "comment_2"

	comments := []*models.Comment{
		{ID: "comment_1", PostID: "post_1", Content: "Корневой", Replies: []*models.Comment{}},
		{ID: "comment_2", PostID: "post_1", ParentID: &parentID1, Content: "Ответ", Replies: []*models.Comment{}},
		{ID: "comment_3", PostID: "post_1", ParentID: &parentID2, Content: "Ответ на ответ", Replies: []*models.Comment{}},
	}

	result := resolver.buildCommentTree(comments)

	// Должен быть только 1 корневой комментарий
	if len(result) != 1 {
		t.Errorf("Ожидали 1 корневой комментарий, получили %d", len(result))
	}

	root := result[0]
	// У корневого комментария должен быть 1 ответ
	if len(root.Replies) != 1 {
		t.Errorf("Ожидали 1 ответ, получили %d", len(root.Replies))
	}

	// Можно добавить проверку второго уровня вложенности
	reply := root.Replies[0]
	if len(reply.Replies) != 1 {
		t.Errorf("Ожидали 1 ответ на ответ, получили %d", len(reply.Replies))
	}
}

// TestBuildCommentTree_BrokenParent - проверяет поведение при "битых" ссылках
func TestBuildCommentTree_BrokenParent(t *testing.T) {
	resolver := &ResolverContext{
		Storage: storage.NewMemoryStorage(),
	}

	// Создаем ID родителя, которого нет в списке комментариев
	brokenParentID := "comment_404"

	comments := []*models.Comment{
		{
			ID:       "comment_1",
			PostID:   "post_1",
			ParentID: &brokenParentID,
			Content:  "Сирота",
			Replies:  []*models.Comment{},
		},
		{
			ID:      "comment_2",
			PostID:  "post_1",
			Content: "Нормальный корневой", 
			Replies: []*models.Comment{},
		},
	}

	result := resolver.buildCommentTree(comments)

	// В текущей реализации: comment_1 теряется, comment_2 становится корневым
	// Ожидаем только 1 корневой комментарий (comment_2)
	if len(result) != 1 {
		t.Errorf("Ожидали 1 корневой комментарий, получили %d", len(result))
	}

	// Проверяем что сохранился именно comment_2, а не comment_1
	if len(result) == 1 && result[0].ID != "comment_2" {
		t.Errorf("Ожидали comment_2, получили %s", result[0].ID)
	}

}
