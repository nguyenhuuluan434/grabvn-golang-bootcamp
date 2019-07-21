package main

import (
	"bufio"
	"fmt"
	"os"
)

//multi operation in one line
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		input := scanner.Text()
		output, err := calculate(input)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(fmt.Sprintf("Value after calculation from input %.f", output))
		}
		fmt.Print("> ")
	}
}
