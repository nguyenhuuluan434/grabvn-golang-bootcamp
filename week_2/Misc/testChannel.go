package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(8)
	var wg sync.WaitGroup
	channel := make(chan int, 10)

	wg.Add(1)
	go func() {
		for i := 0; i <= 50; i++ {
			channel <- i
		}
		fmt.Println("done feed data to channel")
		close(channel)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for value := range channel {
			fmt.Println(value)
			time.Sleep(1 * time.Second)
			fmt.Println("length of channel ", len(channel))
		}
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("done")
}
