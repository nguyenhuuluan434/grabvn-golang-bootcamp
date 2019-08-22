package service

import "grabvn-golang-bootcamp/week_4/theory/dependencyInjection/dao"

type greeter struct {
	database dao.DBService
	lang     string
}

func (g *greeter) Greet() string {
	msg, _ := g.database.FetchMessage(g.lang)
	return "Message is: " + msg
}

func (g *greeter) GreetInDefaultMsg() string {
	msg, _ := g.database.FetchDefaultMessage()
	return "Message is: " + msg
}

type GreetingService interface {
	Greet() string
	GreetInDefaultMsg() string
}

func NewGreeter(dbService dao.DBService, lang string) GreetingService {
	return &greeter{database: dbService, lang: lang}
}




