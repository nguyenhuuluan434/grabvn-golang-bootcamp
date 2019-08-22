package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

func main() {
	greetingMessage := hello("John")
	fmt.Println(aurora.Yellow(greetingMessage))
	greetingMessage = hello("")
	fmt.Println(aurora.Yellow(greetingMessage))

}
