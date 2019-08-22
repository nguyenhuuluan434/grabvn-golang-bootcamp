package service

import (
	"grabvn-golang-bootcamp/week_5/practice/repository"
)

type ServiceImpl struct {
	repository repository.RespositoryService
}

type Service interface {
	CreatePost(request PostRequest) (string, error)
	GetPost(id string) (PostResponse, error)
	DeletePost(id string) (bool, error)
	UpdatePost(id string, request PostRequest) (bool, error)

	CreateComment(postId string, request CommentRequest) (string, error)
	GetComment(id string) (CommentResponse, error)

	GetAllPost() (AllPostWithComments, error)
}

func NewService(daoService repository.RespositoryService) Service {
	return &ServiceImpl{repository: daoService}
}

func (service ServiceImpl) CreatePost(request PostRequest) (id string, err error) {
	requestPostData := repository.RequestPostData{Title: request.Title}
	id, err = service.repository.CreatePost(requestPostData)
	return
}

func (service ServiceImpl) GetPost(id string) (postResponse PostResponse, err error) {
	post, err := service.repository.GetPost(id)
	if err != nil {
		return
	}
	postResponse = PostResponse{Id: post.ID, Title: post.Title}
	return
}

func (service ServiceImpl) DeletePost(id string) (ok bool, err error) {
	ok, err = service.repository.DeletePost(id)
	return
}

func (service ServiceImpl) UpdatePost(id string, request PostRequest) (ok bool, err error) {
	requestPostData := repository.RequestPostData{Title: request.Title}
	ok, err = service.repository.UpdatePost(id, requestPostData)
	return
}

func (service ServiceImpl) CreateComment(postId string, request CommentRequest) (id string, err error) {
	requestCommentData := repository.RequestCommentData{Body: request.Body}
	id, err = service.repository.CreateComment(id, requestCommentData)
	return
}

func (service ServiceImpl) GetComment(id string) (commentResponse CommentResponse, err error) {
	comment, err := service.repository.GetComment(id)
	if err != nil {
		return
	}
	commentResponse = CommentResponse{Id: comment.ID, Body: comment.Body}
	return
}

func (service ServiceImpl) GetAllPost() (postWithComments AllPostWithComments, err error) {
	allPostWithComments,err := service.repository.GetAllPostWithComments()
	if err!=nil {
		return
	}
	if len(allPostWithComments)>0 {
		posts := make([]PostResponse,len(allPostWithComments),len(allPostWithComments))
		postWithComments = AllPostWithComments{}
		for _,post := range allPostWithComments {
			postResponse := PostResponse{Id:post.ID,Title:post.Title}
			if len(post.Comments)>0  {
				commentResponses := make([]CommentResponse,len(post.Comments),len(post.Comments))
				for _,comment :=range post.Comments{
					commentResponses = append(commentResponses, CommentResponse{Id:comment.ID,Body:comment.Body})
				}
				postResponse.Comments= commentResponses
			}
			posts = append(posts, postResponse)
		}
	}
	return
}
