package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//type connectionInfo struct {
//	userName string
//	password string
//	port     int32
//	host     string
//}

type DatabaseRepository interface {
	db() *gorm.DB
	CreateTable([]interface{}) (bool,error)
}

type mySqlConnection struct {
	//connectionInfo *connectionInfo
	connection *gorm.DB
}

func (mysql *mySqlConnection) CreateTable(tables []interface{}) (bool, error) {
	if len(tables) > 0 {
		err := mysql.connection.CreateTable(tables).Error
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return true, nil
}

func (mysql *mySqlConnection) db() *gorm.DB {
	return mysql.connection
}

func NewMySQLRepository(userName string, password string, host string, port int32) (repository DatabaseRepository, err error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/todolist?charset=utf8&parseTime=True", userName, password, host, port))
	if err != nil {
		return nil, err
	}
	repository =  &mySqlConnection{connection: db}
	return
}
