package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)



func main() {

	fmt.Println("Enter input type: <number 1> <operand> <number 2> ")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			continue
		}
		//remove duplicate space in input, strim space begin and end
		regexSpace := regexp.MustCompile(`\s+`)
		input = regexSpace.ReplaceAllString(strings.Trim(input, " "), " ")

		var inputParse = strings.Split(input, " ")

		if len(inputParse) < 3 {
			fmt.Println("Missing parameter!!!, Example 2 + 3 ")
			continue
		}

		if !contains(supportOperations, inputParse[1]) {
			fmt.Printf("Operation %s is not support \n", inputParse[1])
			continue
		}

		firstOperand, err := strconv.Atoi(inputParse[0])
		if err != nil {
			fmt.Println(err)
			continue
		}
		var secondOperand = 0
		if strings.HasSuffix(inputParse[2], "\n") {
			secondOperand, err = strconv.Atoi(strings.ReplaceAll(inputParse[2], "\n", ""))
		} else {
			secondOperand, err = strconv.Atoi(inputParse[2])
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		var result = 0
		switch {
		case inputParse[1] == "+":
			result = firstOperand + secondOperand
			fmt.Printf("%d + %d = %d \n", firstOperand, secondOperand, result)
			continue
		case inputParse[1] == "-":
			result = firstOperand - secondOperand
			fmt.Printf("%d - %d = %d \n", firstOperand, secondOperand, result)
			continue
		case inputParse[1] == "*":
			result = firstOperand * secondOperand
			fmt.Printf("%d * %d = %d \n", firstOperand, secondOperand, result)
			continue
		case inputParse[1] == "/":
			if secondOperand == 0{
				println("Not divide 0")
				continue
			}
			result = firstOperand / secondOperand
			fmt.Printf("%d / %d = %d \n", firstOperand, secondOperand, result)
			continue
		default:
			fmt.Printf("Operation %s is not support \n", inputParse[1])
			continue
		}

	}
}

