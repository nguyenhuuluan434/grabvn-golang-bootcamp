package main

import "fmt"

func main() {
	number := make(chan int)
	go printNumber(number)

	for i := 0; i <= 5; i++ {
		number <- i
	}

	fmt.Scan()
}

func printNumber(message chan int) {
	for msg := range message {
		fmt.Println("got ", msg)
	}
}

//go scheduler
