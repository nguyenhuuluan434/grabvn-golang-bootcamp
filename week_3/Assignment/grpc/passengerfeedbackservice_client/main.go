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

	createPassengerRequest := passengerfeedbackservice.CreatePassengerRequest{Name:"LuanNH"}
	responseCreatePassenger, err := client.AddNewPassenger(ctx, &createPassengerRequest)
	if err != nil {
		log.Fatalf("failed to create passenger: %v", err)
	}
	log.Printf("Received passenger info: %v", responseCreatePassenger)


	from := passengerfeedbackservice.Location{X:100,Y:100}
	to := passengerfeedbackservice.Location{X:120,Y:120}
	createBookingRequest := passengerfeedbackservice.CreateBookingRequest{PassengerId:responseCreatePassenger.Id,From:&from,To:&to}
	for i :=0;i<=3 ;i++ {
		resBooking, err := client.AddBooking(ctx, &createBookingRequest)
		if err != nil {
			log.Fatalf("error occurred when create booking : %v", err)
		}
		log.Printf("Received booking: %v", resBooking)

		updateBookingRequest := passengerfeedbackservice.UpdateBookingRequest{BookingCode: resBooking.Booking.BookingCode}
		resUpdateBooking, err := client.UpdateFinishBooking(ctx, &updateBookingRequest)

		if err != nil {
			log.Fatalf("failed to update booking: %v", err)
		}
		log.Printf("Received update info: %v", resUpdateBooking)

		createFeedBackRequest := passengerfeedbackservice.CreateFeedbackRequest{BookingCode: resBooking.Booking.BookingCode, FeedbackMessage: "Driver need drive be careful"}
		resCreFeedBack, err := client.CreateFeedback(ctx, &createFeedBackRequest)

		if err != nil {
			log.Fatalf("failed to feed back: %v", err)
		}
		log.Printf("Received feedback: %v", resCreFeedBack)

		createFeedBackRequest = passengerfeedbackservice.CreateFeedbackRequest{BookingCode: resBooking.Booking.BookingCode, FeedbackMessage: "Driver need drive be careful 1"}
		resCreFeedBack, err = client.CreateFeedback(ctx, &createFeedBackRequest)
		if err != nil {
			log.Fatalf("failed to create feedback: %v", err)
		}
		if  resCreFeedBack !=nil && resCreFeedBack.ResponseStatus!=nil{
			log.Printf("Error when create another feedback for the same booking %v",resCreFeedBack.ResponseStatus  )
		}else {
			log.Printf("Received create feedback: %v", resCreFeedBack)
		}

		getFeedbackByBookingCode := passengerfeedbackservice.GetFeedbackByBookingCode{BookingCode: resBooking.Booking.BookingCode}

		resGetFeedBackByBookingCode, err := client.GetFeedbackByBookingCode(ctx, &getFeedbackByBookingCode)
		if err != nil {
			log.Fatalf("failed to get feedback by booking code %v", err)
		}
		log.Printf("Received feedback: %v", resGetFeedBackByBookingCode)
	}
	reqGetBookingByPassengerId := passengerfeedbackservice.GetBookingRequestByPassengerId{PassengerId:responseCreatePassenger.Id,PagingRequest:&passengerfeedbackservice.PagingRequest{PageSize:10,PageNumber:1}}
	responseGetBookingByPassengerId,err := client.GetBookingByPassengerId(ctx,&reqGetBookingByPassengerId)
	if err != nil {
		log.Fatalf("failed to get booking by passenger id: %v", err)
	}
	log.Printf("Received booking get by passenger ID: %v", responseGetBookingByPassengerId)

	reqGetFeedbackByPassengerId := passengerfeedbackservice.GetFeedbackByPassengerId{PassengerId:responseCreatePassenger.Id}
	responseGetFeedbackByPassengerId,err := client.GetFeedbackByPassengerId(ctx,&reqGetFeedbackByPassengerId)
		if err != nil {
		log.Fatalf("failed in get feed back by passenger id: %v", err)
	}
	log.Printf("Received feedback by passenger ID: %v", responseGetFeedbackByPassengerId)

	deleteFeedbackByPassengerId:=	passengerfeedbackservice.DeleteFeedbackByPassengerId{PassengerId:responseCreatePassenger.Id}
	responseDeleteFeedbackByPassengerId,err := client.DeleteFeedbackByPassengerId(ctx,&deleteFeedbackByPassengerId )
	if err != nil {
		log.Fatalf("error occurred when delete feed back by passenger id: %v", err)
	}
	log.Printf("received data when call delete feedback by passenger id: %v", responseDeleteFeedbackByPassengerId)


	responseGetFeedbackByPassengerId,err = client.GetFeedbackByPassengerId(ctx,&reqGetFeedbackByPassengerId)
	if err != nil {
		log.Fatalf("failed in get feed back by passenger id: %v", err)
	}
	log.Printf("Received feedback by passenger ID after delete: %v", responseGetFeedbackByPassengerId)


}
