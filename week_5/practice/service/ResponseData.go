package service

type PostResponse struct {
	Id       string            `json:"id"`
	Title    string            `json:"title"`
	Comments []CommentResponse `json:"comments,omitempty"`
}

type CommentResponse struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}

type AllPostWithComments struct {
	Posts []PostResponse `json:"posts"`
}
