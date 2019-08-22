package main

import (
	"fmt"
	"grabvn-golang-bootcamp/week_4/theory/dependencyInjection/dao"
	"grabvn-golang-bootcamp/week_4/theory/dependencyInjection/service"
)

func main() {
	d := dao.NewDB()

	g := service.NewGreeter(d, "en")
	fmt.Println(g.Greet())             // Message is: hello
	fmt.Println(g.GreetInDefaultMsg()) // Message is: default message

	g = service.NewGreeter(d, "es")
	fmt.Println(g.Greet()) // Message is: holla

	g = service.NewGreeter(d, "random")
	fmt.Println(g.Greet()) // Message is: bzzzz
}
