package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Todo struct {
	ID       string
	Title    string
	Complete bool
	Duration int
	CreateAt time.Time
	FinishAt time.Time
}

func (todo *Todo) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New().String())
	return nil
}

const todos string = "todos"

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:1234@tcp(localhost:3306)/todolist?charset=utf8&parseTime=True")
	if err != nil {
		panic("can not connect to database")
		//another stop handle
		//log.Fatal("can not connect to database")
	}
}
func main() {
	defer db.Close()
	db.LogMode(true)

	err1 := db.AutoMigrate(Todo{}).Error

	if err1 != nil {
		log.Fatal("could not create table todo list")
	}

	route := gin.Default()
	route.GET(todos, listTodos)

	route.GET(todos+"/:id", getTodo)

	route.POST(todos+"/:id", updateTodo)
	route.POST(todos, addTodo)
	route.Run(":8787")

}

func listTodos(c *gin.Context) {
	var todos []Todo
	var err error
	err = db.Find(&todos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError{message: "could not get todo list"})
		//c.String(500,"could not get todolist")
		return
	}
	c.JSON(http.StatusOK, todos)
}

func addTodo(c *gin.Context) {
	var arrgument struct {
		Title    string
		Duration int
	}
	err := c.BindJSON(&arrgument)
	if err != nil {
		c.JSON(http.StatusBadRequest, customError{message: "message is invalid, invalid param or bad param"})
		return
	}
	if arrgument.Duration <= 0 {
		c.JSON(http.StatusBadRequest, customError{message: "duration need greater than 0"})
		return
	}

	todo := Todo{Title: arrgument.Title, CreateAt: time.Now(), Complete: false, Duration: arrgument.Duration, FinishAt: time.Now().Add(time.Duration(arrgument.Duration * 3600))}
	//todo := Todo{Title: arrgument.Title}
	err = db.Create(&todo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError{message: "could not save to database"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func getTodo(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, customError{message: "invalid request"})
	}
	var todo Todo

	err := db.Where(&Todo{ID: id}).Find(&todo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError{message: fmt.Sprintf("could not get todo with id %s", id)})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, customError{message: "invalid request"})
	}

	var payload struct {
		Title    string
		Complete bool
	}

	err := c.BindJSON(&payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, customError{message: "bad request"})
	}
	var todo Todo
	var updateError error

	if payload.Complete {
		updateError = db.Model(&todo).Where("id = ?", id).Update(Todo{Title: payload.Title, Complete: payload.Complete, FinishAt: time.Now()}).Error
	} else {
		updateError = db.Model(&todo).Where("id = ?", id).Update(Todo{Title: payload.Title, Complete: payload.Complete}).Error
	}
	if updateError != nil {
		c.JSON(http.StatusInternalServerError, customError{message: fmt.Sprintf("could not update todo have id %s with payload %s", id, payload)})
		return
	}
	c.JSON(http.StatusOK, todo)
}

type customError struct {
	message string
}

func (c *customError) Error() string {
	return c.message
}
