package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	pfs "grabvn-golang-bootcamp/week_3/Assignment/grpc/passengerfeedbackservice"
	"log"
	"net/http"
	"time"
)

type server struct {
	db       *gorm.DB
	logDebug bool
}

func (s server) AddNewPassenger(ctx context.Context, req *pfs.CreatePassengerRequest) (*pfs.PassengerInfoResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	tx := s.db.Begin()
	data := Passenger{ID: uuid.New().String(), Name: req.Name, CreatedAt: time.Now()}
	err := tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		return &pfs.PassengerInfoResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusInternalServerError}}, err
	}
	tx.Commit()
	return &pfs.PassengerInfoResponse{Name: data.Name, Id: data.ID, ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil
}

func (s server) GetPassengerById(ctx context.Context, req *pfs.GetPassengerRequest) (*pfs.PassengerInfoResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	passenger, err := getPassenger(s.db, req.Id)
	if err != nil {
		return &pfs.PassengerInfoResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusInternalServerError}}, err
	}
	if len(passenger.ID) == 0 && len(passenger.Name) == 0 {
		return &pfs.PassengerInfoResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no passenger id %s", req.Id)}}, nil
	}
	return &pfs.PassengerInfoResponse{Id: passenger.ID, Name: passenger.Name, ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil
}

func (s server) AddBooking(cxt context.Context, req *pfs.CreateBookingRequest) (*pfs.CreateBookingResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	passenger, err := getPassenger(s.db, req.PassengerId)
	if err != nil {
		return &pfs.CreateBookingResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no passenger id %s", req.PassengerId)}}, nil
	}
	path, _ := Path{From: req.From, To: req.To}.SerializeToJson()

	data := Booking{ID: uuid.New().String(), PassengerId: passenger.ID, IsDelete: false, IsFinish: false, Path: path, CreatedAt: time.Now(), FinishAt: time.Now()}
	err = s.db.Create(&data).Error
	if err != nil {
		return &pfs.CreateBookingResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusInternalServerError}}, err
	}
	return &pfs.CreateBookingResponse{Booking: &pfs.BookingInfoResponse{IsFinish: data.IsFinish, From: req.From, To: req.To, CreatedAt: &timestamp.Timestamp{Seconds: data.CreatedAt.Unix()}, BookingCode: data.ID}, ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil
}

func (s server) UpdateFinishBooking(ctx context.Context, req *pfs.UpdateBookingRequest) (*pfs.UpdateBookingResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	booking, err := getBooking(s.db, req.BookingCode)
	if err != nil {
		return &pfs.UpdateBookingResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no booking id %s", req.BookingCode)}}, nil
	}
	err = s.db.Model(&booking).Update(Booking{FinishAt: time.Now(), IsFinish: true}).Error
	if err == nil {
		return &pfs.UpdateBookingResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil
	} else {
		return &pfs.UpdateBookingResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusInternalServerError}}, err
	}
}

func (s server) GetBookingByCode(ctx context.Context, req *pfs.GetBookingRequest) (*pfs.GetBookingResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	booking, err := getBooking(s.db, req.BookingCode)
	if err != nil {
		return &pfs.GetBookingResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no booking id %s", req.BookingCode)}}, err
	}
	var from *pfs.Location
	var to *pfs.Location
	path, err := Path{}.DeserializeFromJson(booking.Path)
	if err == nil {
		from = path.From
		to = path.To
	}
	return &pfs.GetBookingResponse{Booking: &pfs.BookingInfoResponse{BookingCode: booking.ID, IsFinish: booking.IsFinish, CreatedAt: &timestamp.Timestamp{Seconds: booking.CreatedAt.Unix()}, FinishAt: &timestamp.Timestamp{Seconds: booking.FinishAt.Unix()}, From: from, To: to}, ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil
}

func (s server) GetBookingByPassengerId(ctx context.Context, req *pfs.GetBookingRequestByPassengerId) (*pfs.PassengerBookingResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	_, err := getPassenger(s.db, req.PassengerId)
	if err != nil {
		return &pfs.PassengerBookingResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no passenger id %s", req.PassengerId)}}, nil
	}
	var rows *sql.Rows
	var rawQuery = "SELECT * FROM bookings WHERE bookings.passenger_id = ? "

	if req.PagingRequest.PageNumber > 0 && req.PagingRequest.PageSize > 0 {
		var start = (req.PagingRequest.PageNumber - 1) * req.PagingRequest.PageSize
		var end = req.PagingRequest.PageNumber * req.PagingRequest.PageSize
		rawQuery = rawQuery + " LIMIT ?,?"
		rows, err = s.db.Raw(rawQuery, req.PassengerId, start, end).Rows()
	} else {
		rows, err = s.db.Raw(rawQuery, req.PassengerId).Rows()
	}
	if err != nil {
		return &pfs.PassengerBookingResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusInternalServerError}}, err
	}
	defer rows.Close()
	var bookings []*pfs.BookingInfoResponse
	for rows.Next() {
		var id, passengerId, path string
		var isFinish, isDelete bool
		var createdAt, finishAt time.Time
		rows.Scan(&id, &passengerId, &path, &isFinish, &isDelete, &createdAt, &finishAt)
		//fmt.Println(rows.Columns())
		var from *pfs.Location
		var to *pfs.Location
		p, err := Path{}.DeserializeFromJson(path)
		if err == nil {
			from = p.From
			to = p.To
		}
		bookings = append(bookings, &pfs.BookingInfoResponse{BookingCode: id, IsFinish: isFinish, CreatedAt: &timestamp.Timestamp{Seconds: createdAt.Unix()}, FinishAt: &timestamp.Timestamp{Seconds: finishAt.Unix()}, From: from, To: to})
	}
	return &pfs.PassengerBookingResponse{Bookings: bookings, ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil

}

func (s server) CreateFeedback(ctx context.Context, req *pfs.CreateFeedbackRequest) (*pfs.CreateFeedbackResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	booking, err := getBooking(s.db, req.BookingCode)
	if err != nil {
		return &pfs.CreateFeedbackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no booking id %s", req.BookingCode)}}, nil
	}
	var feedback Feedback
	err = s.db.Where("booking_id = ?", req.BookingCode).Find(&feedback).Error
	if feedback != (Feedback{}) {
		return &pfs.CreateFeedbackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusConflict, Message: fmt.Sprintf("already add feedback for booking with id %s", req.BookingCode)}}, nil
	}
	data := Feedback{CreatedAt: time.Now(), Message: req.FeedbackMessage, BookingId: booking.ID}
	err = s.db.Create(&data).Error
	if err != nil {
		return &pfs.CreateFeedbackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusInternalServerError}}, err
	}
	return &pfs.CreateFeedbackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil
}

func (s server) GetFeedbackByBookingCode(ctx context.Context, req *pfs.GetFeedbackByBookingCode) (*pfs.GetFeedbackResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	var feedback Feedback
	err := s.db.Where(" booking_id = ? AND is_delete = ?", req.BookingCode, false).Find(&feedback).Error
	if err != nil {
		return &pfs.GetFeedbackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no booking id %s", req.BookingCode)}}, err
	}
	return &pfs.GetFeedbackResponse{FeedBack: &pfs.FeedbackInfoResponse{FeedbackMessage: feedback.Message, BookingCode: feedback.BookingId, CreatedAt: &timestamp.Timestamp{Seconds: feedback.CreatedAt.Unix()}}, ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil

}

func (s server) GetFeedbackByPassengerId(ctx context.Context, req *pfs.GetFeedbackByPassengerId) (*pfs.PassengerFeedbackResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	_, err := getPassenger(s.db, req.PassengerId)
	if err != nil {
		return &pfs.PassengerFeedbackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no passenger with id %s", req.PassengerId)}}, nil
	}

	var rows *sql.Rows
	var rawQuery = " SELECT b.id,f.message,f.created_at "
	rawQuery = rawQuery + "	FROM bookings b "
	rawQuery = rawQuery + " INNER JOIN feedbacks f ON b.passenger_id = ? AND b.id = f.booking_id AND f.is_delete = ? "

	if req.PagingRequest != nil && req.PagingRequest.PageNumber > 0 && req.PagingRequest.PageSize > 0 {
		var start = (req.PagingRequest.PageNumber - 1) * req.PagingRequest.PageSize
		var end = req.PagingRequest.PageNumber * req.PagingRequest.PageSize
		rawQuery = rawQuery + " LIMIT ?,?"
		rows, err = s.db.Raw(rawQuery, req.PassengerId,false, start, end).Rows()
	} else {
		rows, err = s.db.Raw(rawQuery, req.PassengerId,false).Rows()
	}
	if err != nil {
		return &pfs.PassengerFeedbackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusInternalServerError}}, err
	}
	defer rows.Close()
	var feedbacks []*pfs.FeedbackInfoResponse
	for rows.Next() {
		var bookingId, feedBackMessage string
		var createdAt time.Time
		rows.Scan(&bookingId, &feedBackMessage, &createdAt)
		feedbacks = append(feedbacks, &pfs.FeedbackInfoResponse{BookingCode: bookingId, CreatedAt: &timestamp.Timestamp{Seconds: createdAt.Unix()}, FeedbackMessage: feedBackMessage})
	}
	return &pfs.PassengerFeedbackResponse{Feedbacks: feedbacks, ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil
}

func (s server) DeleteFeedbackByPassengerId(ctx context.Context, req *pfs.DeleteFeedbackByPassengerId) (*pfs.DeleteFeedBackResponse, error) {
	if s.logDebug {
		log.Print(req)
	}
	passenger, err := getPassenger(s.db, req.PassengerId)
	if err != nil {
		return &pfs.DeleteFeedBackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusNotFound, Message: fmt.Sprintf("there is no passenger with id %s", req.PassengerId)}}, nil
	}
	var rawQuery = " UPDATE feedbacks f SET f.is_delete = TRUE "
	rawQuery = rawQuery + " WHERE f.booking_id IN ( "
	rawQuery = rawQuery + "		SELECT b.id "
	rawQuery = rawQuery + " 	FROM bookings b "
	rawQuery = rawQuery + " 	WHERE b.passenger_id = ? ) "
	query := s.db.Exec(rawQuery, passenger.ID)
	if query.Error != nil {
		return &pfs.DeleteFeedBackResponse{ResponseStatus: &pfs.ResponseInfo{Code: http.StatusInternalServerError}}, err
	}
	return &pfs.DeleteFeedBackResponse{ TotalFeedBack:query.RowsAffected,ResponseStatus: &pfs.ResponseInfo{Code: http.StatusOK}}, nil

}

func getPassenger(db *gorm.DB, passengerId string) (passenger Passenger, err error) {
	err = db.Where("ID = ?", passengerId).Find(&passenger).Error
	return
}

func getBooking(db *gorm.DB, bookingId string) (booking Booking, err error) {
	err = db.Where("ID = ?", bookingId).Find(&booking).Error
	return
}
