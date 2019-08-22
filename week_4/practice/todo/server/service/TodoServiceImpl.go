package service

import (
	"context"
	"grabvn-golang-bootcamp/week_4/practice/todo/protobuf"
	"grabvn-golang-bootcamp/week_4/practice/todo/server/respository"
	"net/http"
)

type TodoService struct {
	TodoDao respository.ITodoDAO
}

func (t TodoService) CreateTodo(ctx context.Context, req *protobuf.CreateTodoRequest) (*protobuf.CreateTodoResponse, error) {
	id, err := t.TodoDao.Create(req)
	if err != nil {
		return nil, err
	}
	result := protobuf.CreateTodoResponse{Id: id, ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusOK}}
	return &result, nil
}

func (t TodoService) GetTodo(ctx context.Context, req *protobuf.GetTodoRequest) (*protobuf.GetTodoResponse, error) {
	item, err := t.TodoDao.Get(req.Id)
	if err != nil {
		return &protobuf.GetTodoResponse{ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusInternalServerError}}, err
	}
	if item == nil {
		return &protobuf.GetTodoResponse{ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusNotFound}}, nil
	}
	return &protobuf.GetTodoResponse{Item: item, ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusOK}}, nil
}

func (t TodoService) ListTodo(ctx context.Context, req *protobuf.ListTodoRequest) (*protobuf.ListTodoResponse, error) {
	items, err := t.TodoDao.GetList(req.Limit, req.Marker, req.Completed)
	if err != nil {
		return &protobuf.ListTodoResponse{ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusInternalServerError}}, err
	}
	if items == nil {
		return &protobuf.ListTodoResponse{ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusNotFound}}, nil
	}
	return &protobuf.ListTodoResponse{Items: items, ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusOK}}, nil
}

func (t TodoService) UpdateTodo(ctx context.Context, req *protobuf.UpdateTodoRequest) (*protobuf.UpdateTodoReponse, error) {
	err := t.TodoDao.Update(req.Id, req.Item)
	if err != nil {
		return &protobuf.UpdateTodoReponse{ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusInternalServerError}}, err
	}
	return &protobuf.UpdateTodoReponse{ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusOK}}, nil
}

func (t TodoService) DeleteTodo(ctx context.Context, req *protobuf.DeleteTodoRequest) (*protobuf.UpdateTodoReponse, error) {

	err := t.TodoDao.Delete(req.Id)
	if err != nil {
		return &protobuf.UpdateTodoReponse{ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusInternalServerError}}, err
	}
	return &protobuf.UpdateTodoReponse{ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusOK}}, nil
}
