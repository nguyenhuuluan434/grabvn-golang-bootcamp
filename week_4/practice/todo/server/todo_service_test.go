package main

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"grabvn-golang-bootcamp/week_4/practice/todo/protobuf"
	"grabvn-golang-bootcamp/week_4/practice/todo/server/service"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	grpc_runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pact-foundation/pact-go/types"
	"google.golang.org/grpc"
)

var grpcAddress = "localhost:5001"
var httpAddress = "localhost:5002"

type toDoImplMock struct {
}

func (toDoImplMock) Create(item *protobuf.CreateTodoRequest) (string, error) {
	return "id1", nil
}

func (toDoImplMock) Update(id string, item *protobuf.TodoRequestUpdateInfo) error {
	return nil
}

func (toDoImplMock) Get(id string) (*protobuf.Todo, error) {
	return toDos[0], nil
}

func (toDoImplMock) GetList(limit int32, marker string, complete bool) ([]*protobuf.Todo, error) {
	return toDos, nil
}

func (toDoImplMock) Delete(id string) error {
	return nil
}

var toDos []*protobuf.Todo

func startServer() {
	todoRepo := &toDoImplMock{}
	//db := pg.Connect(&pg.Options{
	//	User:                  "postgres",
	//	Password:              "example",
	//	Database:              "todo",
	//	Addr:                  "localhost" + ":" + "5433",
	//	RetryStatementTimeout: true,
	//	MaxRetries:            4,
	//	MinRetryBackoff:       250 * time.Millisecond,
	//})
	s := grpc.NewServer()
	protobuf.RegisterTodoServiceServer(s, service.TodoService{TodoDao: todoRepo})//respository.TodoDAO{DB: db}})

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("can not listen tcp grpcAddress ", grpcAddress, " ", err)
	}

	log.Printf("Serving GRPC at %s.\n", grpcAddress)
	go s.Serve(lis)

	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Couldn't contact grpc server")
	}

	mux := grpc_runtime.NewServeMux()
	err = protobuf.RegisterTodoServiceHandler(context.Background(), mux, conn)
	if err != nil {
		panic("Cannot serve http api")
	}
	log.Printf("Serving http at %s.\n", httpAddress)
	err = http.ListenAndServe(httpAddress, mux)
}

func TestToDoService(t *testing.T) {
	var dir, _ = os.Getwd()
	var pactDir = fmt.Sprintf("%s/../consumer/pacts", dir)
	go startServer()

	pact := &dsl.Pact{
		Consumer: "ToDoConsumer",
		Provider: "ToDoService",
	}
	pact.DisableToolValidityCheck = true
	// Verify the Provider using the locally saved Pact Files
	pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: "http://" + httpAddress,
		PactURLs:        []string{filepath.ToSlash(fmt.Sprintf("%s/todoconsumer-todoservice.json", pactDir))},
		StateHandlers: types.StateHandlers{
			// Setup any state required by the test
			// in this case, we ensure there is a "user" in the system
			"Todo exists": func() error {
				toDos = []*protobuf.Todo{{Id: "id1", Title: "ToDo A"}, {Id: "id2", Title: "ToDo B"}}
				return nil
			},
			"Getting todos existing on Api": func() error {
				toDos = []*protobuf.Todo{{Id: "id1", Title: "ToDo A"}, {Id: "id2", Title: "ToDo B"}}
				return nil
			},

			"Todos existing on Api": func() error {
				toDos = []*protobuf.Todo{{Id: "id1", Title: "ToDo A"}, {Id: "id2", Title: "ToDo B"}}
				return nil
			},
			"A todo with id = id1 exists": func() error {
				toDos = []*protobuf.Todo{{Id: "id1", Title: "ToDo A"}}
				return nil
			},
		},
	})

}
