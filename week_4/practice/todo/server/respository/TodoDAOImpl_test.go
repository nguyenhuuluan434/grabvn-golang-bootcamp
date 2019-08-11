package respository

import (
	"context"
	"fmt"
	protobuf "grabvn-golang-bootcamp/week_4/practice/todo/protobuf"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/suite"
)

type ToDoRepositorySuite struct {
	db *pg.DB
	suite.Suite
	todoDAO TodoDAO
}
type dbLogger struct { }

func (db dbLogger) BeforeQuery(ctx context.Context,q *pg.QueryEvent) (context.Context, error) {
	fmt.Println(q.FormattedQuery())
	return ctx, nil
}

func (db dbLogger) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}

func (s *ToDoRepositorySuite) SetupSuite() {
	// Connect to PostgresQL
	s.db = pg.Connect(&pg.Options{
		User:                  "postgres",
		Password:              "example",
		Database:              "todo",
		Addr:                  "localhost" + ":" + "5433",
		RetryStatementTimeout: true,
		MaxRetries:            4,
		MinRetryBackoff:       250 * time.Millisecond,
	})
	s.db.AddQueryHook(dbLogger{})

	s.db.DropTable(&protobuf.Todo{}, nil)
	// Create Table
	err := s.db.CreateTable(&protobuf.Todo{}, nil)
	if err != nil {
		panic("could not create database")
	}
	s.todoDAO = TodoDAO{DB: s.db}
}

func (s *ToDoRepositorySuite) TearDownSuite() {
	s.db.DropTable(&protobuf.Todo{}, nil)
	s.db.Close()
}

func (s *ToDoRepositorySuite) TestDAO() {
	item := &protobuf.CreateTodoRequest{Item: &protobuf.TodoRequestCreateInfo{Title: "item1", Description: "description for item 1"}}
	id1, err := s.todoDAO.Create(item)
	s.Nil(err)
	item = &protobuf.CreateTodoRequest{Item: &protobuf.TodoRequestCreateInfo{Title: "item2", Description: "description for item 2"}}
	id, err := s.todoDAO.Create(item)

	s.Nil(err)


	newTodo, err := s.todoDAO.Get(id)
	s.Nil(err)
	s.NotNil(newTodo)
	err = s.todoDAO.Update(id,&protobuf.UpdateTodoRequest{Id:id,Item:&protobuf.TodoRequestUpdateInfo{Title:"item 2"}})
	s.Nil(err)

	err = s.todoDAO.Update(id1,&protobuf.UpdateTodoRequest{Id:id1,Item:&protobuf.TodoRequestUpdateInfo{Title:"item klgt",Completed:true}})
	s.Nil(err)

	items ,err :=s.todoDAO.GetList(int32(10),id1,false)
	s.NotEmpty(items)
}

func TestToDoRepository(t *testing.T) {
	suite.Run(t, new(ToDoRepositorySuite))
}
