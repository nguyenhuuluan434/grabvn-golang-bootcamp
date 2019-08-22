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
	t.Run("Test Create ToDo not exist", func(t *testing.T) {
		// Set up our expected interactions.
		reqTitle := "1-1 with manager"
		reqDesc := "discuss about OKRs"
		pact.
			AddInteraction().
			Given("Todo exists").
			UponReceiving("A request to create todo").
			WithRequest(dsl.Request{
				Method:  http.MethodPost,
				Path:    dsl.String("/v1/todo"),
				Body: map[string]interface{}{
					"title":       reqTitle,
					"description": reqDesc,
				},
			}).
			WillRespondWith(dsl.Response{
				Status:  http.StatusOK,
				Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
				Body: dsl.Like(map[string]interface{}{
					"responseStatus": dsl.Like(map[string]interface{}{
						"code": dsl.Like(int32(http.StatusOK)), "message": dsl.Like("none"),
					}),
					"id": dsl.Like("id1"),
				}),
			})

		// Pass in test case. This is the component that makes the external HTTP call
		expected := "id1"
		var test = func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			actual, err := proxy.CreateToDoWithRequest(TodoRequest{Title: reqTitle, Description: reqDesc})
			if err != nil {
				return err
			}

			assert.Equal(t, expected, actual)
			return nil
		}

		// Run the test, verify it did what we expected and capture the contract
		if err := pact.Verify(test); err != nil {
			log.Fatalf("Error on Verify: %v", err)
		}
	})

	t.Run("Test get list ToDo ", func(t *testing.T) {
		pact.AddInteraction().
			Given("Todos existing on Api").
			UponReceiving("A request to list todo").
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
			assert.NotNil(t, resp)
			assert.NotEmpty(t,resp)
			return nil
		}

		if err := pact.Verify(test); err != nil {
			log.Fatalf("Error on Verify: %v", err)
		}
	})

	t.Run("Test Update Todo", func(t *testing.T) {
		id := "id1"
		pact.AddInteraction().
			Given("A todo with id = id1 exists").
			UponReceiving("A request to update todo").
			WithRequest(dsl.Request{
				Method: http.MethodPut,
				Path:   dsl.String(fmt.Sprintf("/v1/todo/%s", id)),
			}).WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(map[string]interface{}{
				"responseStatus": dsl.Like(map[string]interface{}{
					"code": dsl.Like(int32(http.StatusOK)),
				}),
			}),
		})

		expected := int32(200)
		test := func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			actual, err := proxy.UpdateToDo(id, UpdateTodoRequest{Title: "only update title", Completed: true, Description: ""})
			if err != nil {
				return err
			}
			assert.Equal(t, expected, actual)
			return nil
		}

		if err := pact.Verify(test); err != nil {
			log.Printf("Error on Verify: %v", err)
		}
	})

	t.Run("Test Delete Todo", func(t *testing.T) {
		id := "id1"
		pact.AddInteraction().
			Given("A todo with id = id1 exists").
			UponReceiving("A request to delete todo with id1").
			WithRequest(dsl.Request{
				Method: http.MethodDelete,
				Path:   dsl.String(fmt.Sprintf("/v1/todo/%s", id)),
			}).WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(map[string]interface{}{
				//test another way to create data response body
				"responseStatus": map[string]interface{}{
					"code":    int32(http.StatusOK),
				},
			}),
		})
		expectStatus := int32(http.StatusOK)
		test := func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			actualStatus, err := proxy.DeleteToDo(id)
			if err != nil {
				return err
			}
			assert.Equal(t, expectStatus, actualStatus)
			return nil
		}

		if err := pact.Verify(test); err != nil {
			log.Printf("Error on Verify: %v", err)
		}
	})

	t.Run("Test Get Todo", func(t *testing.T) {
		id := "id1"
		title := "ToDo A"
		pact.AddInteraction().
			Given("A todo with id = id1 exists").
			UponReceiving("A request to get todo with id = id1").
			WithRequest(dsl.Request{
				Method: http.MethodGet,
				Path:   dsl.String(fmt.Sprintf("/v1/todo/%s", id)),
			}).WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(map[string]interface{}{
				"item": dsl.Like(map[string]interface{}{"id": dsl.Like(id), "title": dsl.Like(title)}),
				"responseStatus": map[string]interface{}{
					"code":    int32(http.StatusOK),
				},
			}),
		})

		expected := ToDo{Id: id, Title: title}
		test := func() (err error) {
			proxy := ToDoProxy{Host: "localhost", Port: pact.Server.Port}
			actual, err := proxy.GetToDo(id)
			if err != nil {
				return err
			}
			assert.NotNil(t, actual)
			assert.Equal(t, expected, actual)
			return nil
		}

		if err := pact.Verify(test); err != nil {
			log.Printf("Error on Verify: %v", err)
		}
	})
}
