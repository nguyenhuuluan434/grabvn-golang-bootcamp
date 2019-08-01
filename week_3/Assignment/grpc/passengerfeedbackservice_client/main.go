package main

import (
	"context"
	"google.golang.org/grpc"
	passengerfeedbackservice "grabvn-golang-bootcamp/week_3/Assignment/grpc/passengerfeedbackservice"
	"log"
	"time"
)

const (
	port = ":13888"
)
func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect to connect server: %v", err)
	}

	defer conn.Close()
	client := passengerfeedbackservice.NewPassengerFeedbackServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)

	defer cancel()

	createPassengerRequest := passengerfeedbackservice.CreatePassengerRequest{Name:"LuanNH1"}
	resPassenger, err := client.AddNewPassenger(ctx, &createPassengerRequest)
	if err != nil {
		log.Fatalf("failed to receive msg create passenger: %v", err)
	}
	log.Printf("Received passenger: %v", resPassenger)


	from := passengerfeedbackservice.Location{X:100,Y:100}
	to := passengerfeedbackservice.Location{X:120,Y:120}
	createBookingRequest := passengerfeedbackservice.CreateBookingRequest{PassengerId:resPassenger.Id,From:&from,To:&to}
	for i :=0;i<=3 ;i++ {
		resBooking, err := client.AddBooking(ctx, &createBookingRequest)
		if err != nil {
			log.Fatalf("failed to receive msg create booking: %v", err)
		}
		log.Printf("Received booking: %v", resBooking)

		updateBookingRequest := passengerfeedbackservice.UpdateBookingRequest{BookingCode: resBooking.Booking.BookingCode}
		resUpdateBooking, err := client.UpdateFinishBooking(ctx, &updateBookingRequest)

		if err != nil {
			log.Fatalf("failed to receive msg create booking: %v", err)
		}
		log.Printf("Received update booking: %v", resUpdateBooking)

		createFeedBackRequest := passengerfeedbackservice.CreateFeedbackRequest{BookingCode: resBooking.Booking.BookingCode, FeedbackMessage: "Driver need write be careful"}
		resCreFeedBack, err := client.CreateFeedback(ctx, &createFeedBackRequest)

		if err != nil {
			log.Fatalf("failed to receive msg create booking: %v", err)
		}
		log.Printf("Received create feedback: %v", resCreFeedBack)

		getFeedbackByBookingCode := passengerfeedbackservice.GetFeedbackByBookingCode{BookingCode: resBooking.Booking.BookingCode}

		resGetFeedBackByBookingCode, err := client.GetFeedbackByBookingCode(ctx, &getFeedbackByBookingCode)
		if err != nil {
			log.Fatalf("failed to receive msg create booking: %v", err)
		}
		log.Printf("Received feedback: %v", resGetFeedBackByBookingCode)
	}
	reqGetBookingByPassengerId := passengerfeedbackservice.GetBookingRequestByPassengerId{PassengerId:resPassenger.Id,PagingRequest:&passengerfeedbackservice.PagingRequest{PageSize:10,PageNumber:1}}
	responseGetBookingByPassengerId,err := client.GetBookingByPassengerId(ctx,&reqGetBookingByPassengerId)
	if err != nil {
		log.Fatalf("failed to receive msg create booking: %v", err)
	}
	log.Printf("Received booking get by passenger ID: %v", responseGetBookingByPassengerId)

	reqGetFeedbackByPassengerId := passengerfeedbackservice.GetFeedbackByPassengerId{PassengerId:resPassenger.Id}
	responseGetFeedbackByPassengerId,err := client.GetFeedbackByPassengerId(ctx,&reqGetFeedbackByPassengerId)
		if err != nil {
		log.Fatalf("failed to receive msg create booking: %v", err)
	}
	log.Printf("Received feedback by passenger ID: %v", responseGetFeedbackByPassengerId)

}
