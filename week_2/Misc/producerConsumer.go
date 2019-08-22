package main

import (
	"fmt"
	"time"
)

func main() {
	producerConsumer()
	fmt.Scanln()
}

func producerConsumer() {
	channel := make(chan int, 5)

	go func() {
		num := 0
		for {
			num++
			if num == 10 {
				close(channel)
				break
			}
			<-time.After(1 * time.Second)
			channel <- num
		}

	}()

	go func() {
		for {
			select {
			case v, ok := <-channel:
				fmt.Println("received data ", v)
				if !ok {
					return
				}
			}
		}

	}()

	time.Sleep(15 * time.Second)
}
