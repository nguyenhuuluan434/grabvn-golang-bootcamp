package service

import (
	"context"
	"errors"
	"grabvn-golang-bootcamp/week_4/practice/todo/protobuf"
	"grabvn-golang-bootcamp/week_4/practice/todo/server/respository"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTodoServiceCreateTodoOK(t *testing.T) {
	mockToDoRep := &respository.FakeTodoRepo{}
	req := &protobuf.CreateTodoRequest{Item: &protobuf.TodoRequestCreateInfo{Title: "Todo 1", Description: "Description for todo 1"}}

	id := uuid.New().String()
	mockToDoRep.Repo.On("Create", req).Return(id, nil)
	service := TodoService{TodoDao: mockToDoRep}
	result, err := service.CreateTodo(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, id, result.Id)
}

func TestTodoServiceCreateTodoFail(t *testing.T) {
	mockToDoRep := &respository.FakeTodoRepo{}
	req := &protobuf.CreateTodoRequest{Item: &protobuf.TodoRequestCreateInfo{Title: "Todo 1", Description: "Description for todo 1"}}
	mockToDoRep.Repo.On("Create", req).Return(nil, errors.New("duplicate"))
	service := TodoService{TodoDao: mockToDoRep}
	_, err := service.CreateTodo(context.Background(), req)
	assert.Nil(t, err)
}

func TestTodoServiceGetTodoOK(t *testing.T) {
	mockToDoRep := &respository.FakeTodoRepo{}
	id := uuid.New().String()
	todo := &protobuf.Todo{Id: id}
	req := &protobuf.GetTodoRequest{Id: id}
	mockToDoRep.Repo.On("Get", req.Id).Return(todo, nil)
	expected := protobuf.GetTodoResponse{Item: todo, ResponseStatus: &protobuf.ResponseStatus{Code: http.StatusOK}}
	service := TodoService{TodoDao: mockToDoRep}
	result, err := service.GetTodo(context.Background(), req)
	assert.Nil(t, err)
	assert.NotNil(t, expected, result)
}

func TestTodoServiceGetFail(t *testing.T) {
	mockToDoRep := &respository.FakeTodoRepo{}
	id := uuid.New().String()
	req := &protobuf.GetTodoRequest{Id: id}
	mockToDoRep.Repo.On("Get", req.Id).Return(nil, errors.New("not found"))
	service := TodoService{TodoDao: mockToDoRep}
	_, err := service.GetTodo(context.Background(), req)
	assert.NotNil(t, err)
}

func TestTodoServiceUpdateTodoOK(t *testing.T) {
	mockToDoRep := &respository.FakeTodoRepo{}
	id := uuid.New().String()

	req := &protobuf.UpdateTodoRequest{Id: id, Item: &protobuf.TodoRequestUpdateInfo{Title: "updated"}}
	mockToDoRep.Repo.On("Update", req.Id, req.Item).Return(nil)
	service := TodoService{TodoDao: mockToDoRep}
	result, err := service.UpdateTodo(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, int32(http.StatusOK), result.ResponseStatus.Code)
}

func TestTodoServiceListTodoOK(t *testing.T) {
	mockToDoRep := &respository.FakeTodoRepo{}
	items := []*protobuf.Todo{
		{Id: uuid.New().String(), Completed: true, Title: "1234"},
	}
	mockToDoRep.Repo.On("GetList", int32(10), "", false).Return(items, nil)
	service := TodoService{TodoDao: mockToDoRep}
	req := &protobuf.ListTodoRequest{Limit: int32(10), Marker: "", Completed: false}
	result, err := service.ListTodo(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, int32(http.StatusOK), result.ResponseStatus.Code)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Items)
	assert.NotEmpty(t, result.Items)
}

func TestTodoServiceDeleteTodo(t *testing.T) {
	mockToDoRep := &respository.FakeTodoRepo{}
	id := uuid.New().String()
	mockToDoRep.Repo.On("Delete", id).Return( nil)
	service := TodoService{TodoDao: mockToDoRep}
	req :=&protobuf.DeleteTodoRequest{Id:id}
	result, err := service.DeleteTodo(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, int32(http.StatusOK), result.ResponseStatus.Code)
	assert.NotNil(t, result)
}
