package service

import (
	"context"
	"grabvn-golang-bootcamp/week_4/practice/todo/protobuf"
)

type ITodoService interface {
	CreateTodo(ctx context.Context, req *protobuf.CreateTodoRequest) (*protobuf.CreateTodoResponse, error)
	GetTodo(ctx context.Context, req *protobuf.GetTodoRequest) (*protobuf.GetTodoResponse, error)
	ListTodo(ctx context.Context, req *protobuf.ListTodoRequest) (*protobuf.ListTodoResponse, error)
	UpdateTodo(ctx context.Context, req *protobuf.UpdateTodoRequest) (*protobuf.UpdateTodoReponse, error)
	DeleteTodo(ctx context.Context, req *protobuf.DeleteTodoRequest) (*protobuf.UpdateTodoReponse, error)
}
