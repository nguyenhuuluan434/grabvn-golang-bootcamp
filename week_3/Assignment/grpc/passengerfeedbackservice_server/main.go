package main

import (
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	passengerfeedbackservice "grabvn-golang-bootcamp/week_3/Assignment/grpc/passengerfeedbackservice"
	"log"
	"net"
)

const (
	port = ":13888"
	logDebug  = true
)
var db *gorm.DB

//var dbService *DbService


func init() {
	var err error
	db, err = gorm.Open("mysql", "root:1234@tcp(localhost:3306)/passenger?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatal("can not connect to database %v", err)
	}
	db.LogMode(logDebug)
	err = db.AutoMigrate(Passenger{}, Booking{}).AddForeignKey("passenger_id", "passengers(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatal("could not create table Passenger,Booking %v", err)
	}
	err = db.AutoMigrate(Feedback{}).AddForeignKey("booking_id", "bookings(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatal("could not create table FeedBack %v", err)
	}
	//passengerDbService := PassengerDbServiceImpl{db:db}
	//dbService :=DbService{PassengerDbService:&passengerDbService}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := server{db: db,logDebug:logDebug}
	passengerfeedbackservice.RegisterPassengerFeedbackServiceServer(s, &server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
