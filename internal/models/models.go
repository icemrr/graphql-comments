package models


type Post struct {
	ID       string     `json:"id"`       
	Title    string     `json:"title"`    
	Content  string     `json:"content"`  
	Comments []*Comment `json:"comments"` 
}


type Comment struct {
	ID       string     `json:"id"`       
	PostID   string     `json:"postId"`  
	ParentID *string    `json:"parentId"` 
	Content  string     `json:"content"` 
	Replies  []*Comment `json:"replies"`  
}


type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}


type CreateCommentInput struct {
	PostID   string  `json:"postId"`   
	ParentID *string `json:"parentId"` 
	Content  string  `json:"content"`  
}
