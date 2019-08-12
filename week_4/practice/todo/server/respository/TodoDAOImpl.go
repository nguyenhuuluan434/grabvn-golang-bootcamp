package respository

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"grabvn-golang-bootcamp/week_4/practice/todo/protobuf"
	"time"
)

type TodoDAO struct {
	DB *pg.DB
}

func (t TodoDAO) Create(req *protobuf.CreateTodoRequest) (id string, err error) {
	id = uuid.New().String()
	entity := protobuf.Todo{Id: id, Title: req.Item.Title, Description: req.Item.Description, Completed: false, Deleted: false, CreatedAt: &timestamp.Timestamp{Seconds: time.Now().Unix()}, UpdatedAt: &timestamp.Timestamp{Seconds: time.Now().Unix()}}
	err = t.DB.Insert(&entity)
	return
}

func (t TodoDAO) Get(id string) (item *protobuf.Todo, err error) {
	item = &protobuf.Todo{}
	err = t.DB.Model(item).Where("id = ?", id).First()
	return
}

func (t TodoDAO) Update(id string, item *protobuf.TodoRequestUpdateInfo) (err error) {
	todoInDB, err := t.Get(id)
	if err != nil {
		return err
	}
	if todoInDB.Completed != item.Completed {
		todoInDB.Completed = item.Completed
	}
	if item.Description != todoInDB.Description && len(item.Description) > 0 {
		todoInDB.Description = item.Description
	}
	if todoInDB.Title != item.Title && len(todoInDB.Title) > 0 {
		todoInDB.Title = item.Title
	}
	todoInDB.UpdatedAt = &timestamp.Timestamp{Seconds: time.Now().Unix()}
	err = t.DB.Update(todoInDB)
	return
}

func (t TodoDAO) GetList(limit int32, marker string, complete bool) ([]*protobuf.Todo, error) {
	todoMarker, err := t.Get(marker)
	if err != nil || todoMarker == nil {
		return nil, errors.New("invalid marker")
	}
	var items []*protobuf.Todo
	query := t.DB.Model(&items).Order("created_at ASC").Where("CAST (created_at -> 'seconds' AS INTEGER )  >= ?", timestamp.Timestamp{Seconds: todoMarker.CreatedAt.Seconds}.Seconds)
	if limit > 0 {
		query.Limit(int(limit))
	}
	query.Where("completed = ?", complete)
	err = query.Select()
	return items, err
}

func (t TodoDAO) Delete(id string) (err error) {
	err = t.DB.Delete(&protobuf.Todo{Id: id})
	return
}
