package main

import (
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	customerfeedback "grabvn-golang-bootcamp/week_3/Assignment/grpc/passengerfeedbackservice"
	"time"
)

type Passenger struct {
	ID        string `gorm:"type:varchar(50);unique_index"`
	Name      string
	CreatedAt time.Time
	Booking   []Booking //`gorm:"foreignkey:PassengerId"`
}

type Booking struct {
	ID          string `gorm:"type:varchar(50);unique_index"`
	PassengerId string
	//present from to in json
	Path      string
	IsFinish  bool
	IsDelete  bool
	Feedback  Feedback
	CreatedAt time.Time
	FinishAt  time.Time
}

type Feedback struct {
	BookingId string `gorm:"type:varchar(50);unique_index"`
	Message   string `gorm:"type:varchar(2048)"`
	CreatedAt time.Time
	IsDelete  bool
}

type CustomerFeedBack struct {
	CustomerId      string
	CustomerName    string
	BookingId       string
	FeedbackMesaage string
}

type Path struct {
	From *customerfeedback.Location
	To   *customerfeedback.Location
}

func (p Path) SerializeToJson() (string, error) {
	data, err := json.Marshal(p)
	return string(data), err
}

func (p Path) DeserializeFromJson(input string) (path Path, err error) {
	byteData := []byte(input)
	err = json.Unmarshal(byteData, &path)
	return
}

//type DbService struct {
//	PassengerDbService PassengerDbService
//}
//
//type PassengerDbService interface {
//	addPassenger(name string) (Passenger, error)
//}
//
//type PassengerDbServiceImpl struct {
//	db *gorm.DB
//}
//
//func (p PassengerDbServiceImpl) addPassenger(name string) (Passenger, error) {
//	passenger := Passenger{ID: uuid.New().String(), Name: name, CreatedAt: time.Time{}}
//	err := p.db.Create(&passenger).Error
//	if err != nil {
//		return Passenger{}, err
//	}
//	return passenger, nil
//}

//var db *gorm.DB
//
//func init() {
//	var err error
//	db, err = gorm.Open("mysql", "root:1234@tcp(localhost:3306)/passenger?charset=utf8&parseTime=True")
//	if err != nil {
//		panic("can not connect to database")
//		//another stop handle
//		//log.Fatal("can not connect to database")
//	}
//}
//func main() {
//	defer db.Close()
//	defer db.Close()
//	defer db.Close()
//	db.LogMode(true)
//
//	err1 := db.AutoMigrate(Passenger{},Booking{}).AddForeignKey("passenger_id","passengers(id)","RESTRICT", "RESTRICT").Error
//
//	if err1 != nil {
//		log.Fatal("could not create table Passenger,Booking ")
//	}
//	err1 = db.AutoMigrate(Feedback{}).AddForeignKey("booking_id","bookings(id)","RESTRICT", "RESTRICT").Error
//	if err1 != nil {
//		log.Fatal("could not create table FeedBack")
//	}
//}
