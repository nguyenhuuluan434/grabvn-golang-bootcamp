package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//doTask()
	raceConditionExample()
	fixRaceCondtionExample()

}

func doTask() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("do task 1")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("do task 2")
	}()
	fmt.Println("Dispatched task")
	wg.Wait()
	fmt.Println("Done tasks")
}

func raceConditionExample() {

	n := 0
	go func() {
		for i := 0; i < 100000; i++ {
			n++
		}
	}()
	go func() {
		for i := 0; i < 100000; i++ {
			n++
		}
	}()
	time.Sleep(5 * time.Second)
	fmt.Printf("%d\n", n)
}

func fixRaceCondtionExample() {
	n := 0

	lock := &sync.Mutex{}
	go func() {
		for i := 0; i < 100000; i++ {
			lock.Lock()
			n++
			lock.Unlock()
		}
	}()
	go func() {
		for i := 0; i < 100000; i++ {
			lock.Lock()
			n++
			lock.Unlock()
		}
	}()
	time.Sleep(5 * time.Second)
	fmt.Println(n)

}
