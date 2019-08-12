package consumer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type ToDo struct {
	Id          string               `json:"id,omitempty" `
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	Completed   bool                 `json:"completed,omitempty"`
	CreatedAt   *timestamp.Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *timestamp.Timestamp `json:"updated_at,omitempty"`
}

type TodoRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed,omitempty"`
}

type ToDoProxy struct {
	Host string
	Port int
}

type ToDoList struct {
	Items []ToDo `json:"items"`
}

func (p *ToDoProxy) getApi() string {
	return fmt.Sprintf("http://%s:%d/v1", p.Host, p.Port)
}

type ResponseStatus struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (p *ToDoProxy) CreateToDoWithRequest(todo TodoRequest) (id string,err error) {
	toDoBytes, _ := json.Marshal(todo)

	u := fmt.Sprintf("%s/todo", p.getApi())
	req, _ := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(toDoBytes))
	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	createToDoRes := struct {
		ID             string `json:"id"`
		ResponseStatus ResponseStatus `json:"responseStatus"`
	}{}

	json.NewDecoder(res.Body).Decode(&createToDoRes)
	if createToDoRes.ResponseStatus.Code != int32(http.StatusOK) {
		err = errors.New(createToDoRes.ResponseStatus.Message)
	}
	return createToDoRes.ID, err
}

func (p *ToDoProxy) GetListToDo(limit int32, maker string, completed bool) (*ToDoList, error) {
	url := fmt.Sprintf("%s/todo?limit=%d&completed=%t&marker=%s", p.getApi(), limit, completed, maker)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	toDos := ToDoList{}
	json.NewDecoder(res.Body).Decode(&toDos)
	return &toDos, nil
}

func (p *ToDoProxy) UpdateToDo(id string, todo UpdateTodoRequest) (statusCode int32, err error) {
	toDoBodyReqInByte, _ := json.Marshal(todo)
	url := fmt.Sprintf("%s/todo/%s", p.getApi(), id)
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(toDoBodyReqInByte))

	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	updateToDoRes := struct {
		ResponseStatus ResponseStatus `json:"responseStatus"`
	}{}
	json.NewDecoder(res.Body).Decode(&updateToDoRes)
	if updateToDoRes.ResponseStatus.Code != int32(http.StatusOK) {
		err = errors.New(updateToDoRes.ResponseStatus.Message)
	}
	return updateToDoRes.ResponseStatus.Code, err
}

func (p *ToDoProxy) DeleteToDo(id string) (statusCode int32, err error) {

	url := fmt.Sprintf("%s/todo/%s", p.getApi(), id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	deleteToDoRes := struct {
		ResponseStatus ResponseStatus `json:"responseStatus"`
	}{}
	json.NewDecoder(res.Body).Decode(&deleteToDoRes)
	if deleteToDoRes.ResponseStatus.Code != int32(http.StatusOK) {
		err = errors.New(deleteToDoRes.ResponseStatus.Message)
	}
	return deleteToDoRes.ResponseStatus.Code, err
}

func (p *ToDoProxy) GetToDo(id string) (todo ToDo, err error) {

	url := fmt.Sprintf("%s/todo/%s", p.getApi(), id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	getToDoRes := struct {
		todo           ToDo `json:"todo"`
		ResponseStatus ResponseStatus `json:"responseStatus"`
	}{}
	json.NewDecoder(res.Body).Decode(&getToDoRes)
	if getToDoRes.ResponseStatus.Code !=int32(http.StatusOK) {
		err = errors.New(getToDoRes.ResponseStatus.Message)
	}
	return getToDoRes.todo, err
}
