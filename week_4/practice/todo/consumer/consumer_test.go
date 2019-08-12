package consumer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestToDoProxy(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "ToDoConsumer",
		Provider: "ToDoService",
		Host:     "localhost",
	}
	defer pact.Teardown()
	env := os.Getenv("PATH")
	_ = os.Setenv("PATH", env+":/usr/bin/pact/bin/")

	pact.DisableToolValidityCheck = true
	t.Run("TestCreateToDo", func(t *testing.T) {

		// Set up our expected interactions.
		pact.
			AddInteraction().
			//Given("UserA is existing").
			UponReceiving("A request to create todo").
			WithRequest(dsl.Request{
				Method:  http.MethodPost,
				Path:    dsl.String("/v1/todo"),
				Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
				Body: map[string]interface{}{
					"title":       "1-1 with manager",
					"description": "discuss about OKRs",
				},
			}).
			WillRespondWith(dsl.Response{
				Status:  http.StatusOK,
				Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
				Body: dsl.Like(map[string]interface{}{
					"responseStatus": dsl.Like(map[string]interface{}{
						"code": dsl.Like(int32(http.StatusOK)), "message": dsl.Like("klgt"),
					}),
					"id": dsl.Like("id1"),
				}),
			})

		// Pass in test case. This is the component that makes the external HTTP call
		var test = func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			id, err := proxy.CreateToDoWithRequest(TodoRequest{Title: "1-1 with manager", Description: "discuss about OKRs"})
			if err != nil {
				return err
			}
			assert.Equal(t, "id1", id)
			return nil
		}

		// Run the test, verify it did what we expected and capture the contract
		if err := pact.Verify(test); err != nil {
			log.Fatalf("Error on Verify: %v", err)
		}
	})

	t.Run("TestGetToDoList", func(t *testing.T) {
		pact.AddInteraction().Given("There are todos").UponReceiving("A request to list todo").
			WithRequest(dsl.Request{
				Method: http.MethodGet,
				Path:   dsl.String("/v1/todo"),
				Query:  dsl.MapMatcher{"limit": dsl.String("10"), "completed": dsl.String("true"), "marker": dsl.String("id1")},
			}).WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(map[string]interface{}{
				"items": dsl.Like([]interface{}{
					dsl.Like(map[string]interface{}{"id": dsl.Like("id1"), "title": dsl.Like("ToDo A")}),
					dsl.Like(map[string]interface{}{"id": dsl.Like("id2"), "title": dsl.Like("ToDo B")}),
				}),
			}),
		})

		test := func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			resp, err := proxy.GetListToDo(10, "id1", true)

			if err != nil {
				return err
			}
			fmt.Println(resp)
			return nil
		}

		if err := pact.Verify(test); err != nil {
			log.Fatalf("Error on Verify: %v", err)
		}
	})

	t.Run("TestUpdateTodo", func(t *testing.T) {
		id := "id1"
		pact.AddInteraction().Given("Update todo with id = id1").UponReceiving("A request to list todo").
			WithRequest(dsl.Request{
				Method: http.MethodPut,
				Path:   dsl.String(fmt.Sprintf("/v1/todo/%s", id)),
			}).WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(map[string]interface{}{
				"responseStatus": dsl.Like(map[string]interface{}{
					"code": dsl.Like(int32(http.StatusOK)), "message": dsl.Like(""),
				}),
			}),
		})

		test := func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			resp, err := proxy.UpdateToDo(id, UpdateTodoRequest{Title: "only update title", Completed: true, Description: ""})
			if err != nil {
				return err
			}
			fmt.Println(resp)
			return nil
		}

		if err := pact.Verify(test); err != nil {
			log.Printf("Error on Verify: %v", err)
		}
	})

	t.Run("TestDeleteTodo", func(t *testing.T) {
		id := "id1"
		pact.AddInteraction().Given("Delete todo with id = id1").UponReceiving("A request to list todo").
			WithRequest(dsl.Request{
				Method: http.MethodDelete,
				Path:   dsl.String(fmt.Sprintf("/v1/todo/%s", id)),
			}).WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(map[string]interface{}{
				//test another way to create data response body
				"responseStatus":map[string]interface{}{
					"name":int32(http.StatusOK) ,
					"message": "",
				},
			}),
		})

		test := func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			resp, err := proxy.DeleteToDo(id)
			if err != nil {
				return err
			}
			fmt.Println(resp)
			return nil
		}

		if err := pact.Verify(test); err != nil {
			log.Printf("Error on Verify: %v", err)
		}
	})

	t.Run("TestGetTodo", func(t *testing.T) {
		id := "id1"
		pact.AddInteraction().Given("Get todo with id = id1").UponReceiving("A request to list todo").
			WithRequest(dsl.Request{
				Method: http.MethodGet,
				Path:   dsl.String(fmt.Sprintf("/v1/todo/%s", id)),
			}).WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(map[string]interface{}{
				"todo": dsl.Like(map[string]interface{}{"id": dsl.Like("id1"), "title": dsl.Like("ToDo A")}),
				"responseStatus": map[string]interface{}{
					"name":    int32(http.StatusOK),
					"message": "",
				},
			}),
		})

		test := func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			resp, err := proxy.GetToDo(id)
			if err != nil {
				return err
			}
			fmt.Println(resp)
			return nil
		}

		if err := pact.Verify(test); err != nil {
			log.Printf("Error on Verify: %v", err)
		}
	})
}
