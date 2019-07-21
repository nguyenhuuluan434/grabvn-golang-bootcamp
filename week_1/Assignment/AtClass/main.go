package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		input := scanner.Text()
		output, err := eval(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(output)
		fmt.Print("> ")
	}
}
