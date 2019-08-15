//+build wireinject
package main

import (
	"bytes"
	"fmt"
	"github.com/google/wire"
)

func main() {
	logger := &LoggerDI{}
	describe
	service := NewCollectServiceDI(logger,NewHttpClientDI(logger))
	res := service.GetAll("https://www.google.com.vn", "https://vng.com.vn")
	fmt.Println(res)
}

type LoggerDI struct {
}

func (logger *LoggerDI) Log(message string) {
	fmt.Println(message)
}

type HttpClientDI struct {
	*LoggerDI
}

func (client *HttpClientDI) Get(url string) string {
	client.Log("Getting " + url)
	return "my response from " + url
}

func NewHttpClientDI(logger *LoggerDI) *HttpClientDI {
	return &HttpClientDI{logger}
}

type CollectServiceDI struct {
	*LoggerDI
	*HttpClientDI
}

func (service *CollectServiceDI) GetAll(urls ...string) string {
	service.Log("Begin collect from urls ...")
	var result bytes.Buffer
	for _, url := range urls {
		result.WriteString(service.Get(url))
	}
	return result.String()
}
func NewCollectServiceDI(logger *LoggerDI, client *HttpClientDI) *CollectServiceDI {
	return &CollectServiceDI{logger, client}
}

func CreateCollectServiceDI() *CollectServiceDI {
	//logger :=  &Logger{}
	//client := NewHttpClientDI(logger)
	//return &CollectService{ logger, client}
	panic(	wire.Bind(HttpClientDI{},NewCollectServiceDI))
	panic(wire.Bind(LoggerDI{},NewCollectServiceDI))
}
