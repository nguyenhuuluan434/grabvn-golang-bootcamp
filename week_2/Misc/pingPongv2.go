package main

import (
	"fmt"
	"time"
)

type Ball struct {
	hits int
}

func main() {
	//create channel un-buffer
	table := make(chan *Ball)

	go player("ping", table)
	go player("pong", table)

	table <- new(Ball)
	time.Sleep(1 * time.Second)
	//get latest element out channel
	<-table
}

func player(name string, table chan *Ball) {
	for ball := range table {
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
