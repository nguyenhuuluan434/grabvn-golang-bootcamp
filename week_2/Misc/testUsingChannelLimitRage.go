package main

import (
	"fmt"
	"time"
)

func main() {
	nums := gen()
	pool := make(chan int, 10)
	for {
		var v int
		var ok bool
		select {
		case v, ok = <-nums:
			if !ok {
				break
			}
			pool <- 1
			if len(pool) <= 10 {
				fmt.Println("len ", len(pool))
				go do(v, pool)
			}
		}
		if !ok {
			break
		}
	}
	fmt.Scanln()
}

func gen() <-chan int {
	c := make(chan int, 10)
	go func() {
		defer close(c)
		for i := 0; i <= 100; i++ {
			c <- i
		}
	}()
	return c
}

func do(i int, pool chan int) {
	time.Sleep(200 * time.Millisecond)
	fmt.Println(i)
	<-pool
}
