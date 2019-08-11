package respository

import "grabvn-golang-bootcamp/week_4/practice/todo/protobuf"

type ITodoDAO interface {
	Create(item *protobuf.CreateTodoRequest) (string, error)
	Update(id string, item *protobuf.TodoRequestUpdateInfo) error
	Get(id string) (*protobuf.Todo, error)
	GetList(limit int32, marker string, complete bool) ([]*protobuf.Todo, error)
	Delete(id string) error
}
