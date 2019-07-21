package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Todo struct {
	ID       int
	Title    string
	Complete bool
	CreateAt time.Time
	FinishAt time.Time
}

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

	err := db.AutoMigrate(Todo{}).Error

	if err != nil {
		log.Fatal("could not create table todolist")
	}

	route := gin.Default()
	route.GET("/todo", listTodos)

	route.POST("/todo", addTodo)
	route.Run(":8787")

}

func listTodos(c *gin.Context) {
	var todos []Todo
	var err error
	err = db.Find(&todos).Error
	if err != nil {
		c.JSON(500, customError{message: "could not get todolist", err: err})
		//c.String(500,"could not get todolist")
		return
	}
	c.JSON(200, todos)
}

func addTodo(c *gin.Context) {
	var arrgument struct {
		Title string
	}
	err := c.BindJSON(&arrgument)
	if err != nil {
		c.JSON(400, customError{message: "Message is invalid, invalid param or bad param", err: err})
		return
	}

	todo := Todo{Title: arrgument.Title,CreateAt:time.Now(),FinishAt:time.Now()}
	//todo := Todo{Title: arrgument.Title}
	err = db.Create(&todo).Error
	if err != nil {
		c.JSON(500, customError{message: "could not save to database"})
		return
	}
	c.JSON(200, todo)
}

type customError struct {
	message string
	err     error
}

func (c *customError) Error() string {
	return c.message
}
