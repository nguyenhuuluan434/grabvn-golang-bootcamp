package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	passengerfeedbackservice "grabvn-golang-bootcamp/week_3/Assignment/grpc/passengerfeedbackservice"
	"log"
	"net/http"
	"time"
)

const (
	grpcPort    = ":13888"
	restApiPort = ":8787"
)

var client passengerfeedbackservice.PassengerFeedbackServiceClient
var ctx context.Context
var cancel context.CancelFunc
var conn grpc.ClientConn

func init() {
	conn, err := grpc.Dial("localhost"+grpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect to connect server: %v", err)
	}

	client = passengerfeedbackservice.NewPassengerFeedbackServiceClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*2)
}

func main() {
	route := gin.Default()
	route.POST("passengers", createPassenger)
	route.Run(restApiPort)
	defer cancel()
	defer conn.Close()
}

func createPassenger(c *gin.Context) {
	var arrgument struct {
		Name string
	}
	err := c.BindJSON(&arrgument)
	log.Println(arrgument)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("bad request"))
		return
	}
	createPassengerRequest := passengerfeedbackservice.CreatePassengerRequest{Name: "LuanNH1"}
	resPassenger, err := client.AddNewPassenger(ctx, &createPassengerRequest)
	if err != nil {
		log.Println("failed to receive msg create passenger: %v", err)
		c.JSON(http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	var dataResponse struct {
		Name string
		Id   string
	}
	dataResponse.Id = resPassenger.Id
	dataResponse.Name = resPassenger.Name
	c.JSON(http.StatusOK, dataResponse)
	return
}
