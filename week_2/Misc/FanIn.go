package main

import (
	"fmt"
)

func main() {
	var data []int
	for i := 0; i < 10; i++ {
		data = append(data, i)
	}
	orgData := generateData(data)
	multiByOne := make(chan int, 10)
	multiByTwo := make(chan int, 10)
	multiByThree := make(chan int, 10)

	go spread(orgData, multiByOne, multiByTwo, multiByThree)

	go processChannel(multiByOne, 1, multiplier)

	go processChannel(multiByTwo, 2, multiplier)

	go processChannel(multiByThree, 3, multiplier)

	fmt.Scanln()
}
func generateData(nums []int) chan int {
	result := make(chan int, 30)
	for value := range nums {
		result <- value
	}
	return result
}

func multiplier(multiTime int, orgInput int) int {
	return multiTime * orgInput
}

func spread(main, a, b, c chan int) {
	for value := range main {
		a <- value
		b <- value
		c <- value
	}
}

func printData(input chan int) {
	for value := range input {
		fmt.Println(value)
	}
}

func processChannel(input chan int, multiTime int, x func(int, int) int) {
	for dataInput := range input {
		result := x(multiTime, dataInput)
		fmt.Println("multi by", multiTime, result)
	}
}

func FanInABC(a, b, c chan int) {
	var total = 0
	select {
	case v := <-a:
		total += v
	case v := <-b:
		total += v
	case v := <-c:
		total += v
	}
	fmt.Println(total)
}
