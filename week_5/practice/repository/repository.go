package repository

type RespositoryService interface {
	CreatePost(req RequestPostData) (string, error)
	UpdatePost(postId string, req RequestPostData) (bool, error)
	GetPost(postId string) (Post, error)
	DeletePost(postId string) (bool, error)

	CreateComment(postId string, req RequestCommentData) (string, error)
	GetComment(commentId string) (Comment, error)
	DeleteComment(commentId string) (bool, error)

	GetAllPostWithComments() ([]PostWithComments, error)
}


type RequestPostData struct {
	Title   string `json:"title"`
}

type RequestCommentData struct {
	Body string `json:"body"`
}

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
}
type Comment struct {
	ID     string `json:"id"`
	Body   string `json:"body"`
	PostID string `json:"postId"`
}
type PostWithCommentsResponse struct {
	Posts []PostWithComments `json:"posts"`
}
type PostWithComments struct {
	ID       string    `json:"id"`
	Title    string    `json:"string"`
	Comments []Comment `json:"comments,omitempty"`
}