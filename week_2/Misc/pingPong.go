package main

import "fmt"

func main() {
	var ping = make(chan int)
	var pong = make(chan int)

	go pinger(ping, pong)
	go ponger(pong)
	ping <- 1
	fmt.Scanln()
}

func pinger(ping chan int, pong chan int) {
	for {
		<-ping

		fmt.Println("ping")
		pong <- 1
	}
}

func ponger(pong chan int) {
	for range pong {
		fmt.Println("Pong")
	}
}
