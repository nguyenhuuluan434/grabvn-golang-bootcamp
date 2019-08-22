package service

type PostRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CommentRequest struct {
	Body   string `json:"body"`
}