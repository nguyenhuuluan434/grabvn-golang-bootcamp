package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var supportOperations = []string{"+", "-", "*", "/"}

func eval(input string) (result string, err error) {
	a, b, op, err := parse(input)
	if err != nil {
		return
	}
	if !contains(supportOperations, op) {
		return result, errors.New(fmt.Sprintf("Not support op: %s", op))
	}
	switch op {
	case "+":
		return fmt.Sprintf(" %.f + %.f = %.f", a, b, a+b), nil
	case "-":
		return fmt.Sprintf(" %.f - %.f = %.f", a, b, a-b), nil
	case "*":
		return fmt.Sprintf(" %.f * %.f = %.f", a, b, a*b), nil
	case "/":
		if b == 0 {
			return "", errors.New("Not divide by 0")
		}
		return fmt.Sprintf(" %.f / %.f = %.f", a, b, a/b), nil
	default:
		return "", errors.New(fmt.Sprintf("Not support op: %s", op))
	}
}

func parse(input string) (firstOperand, secondOperand float64, op string, err error) {
	//replace multi space with one space to split
	regexSpace, err := regexp.Compile(`\s+`)
	if err != nil {
		return
	}
	input = regexSpace.ReplaceAllString(strings.Trim(input, " "), " ")

	var inputParse = strings.Split(input, " ")

	if len(inputParse) < 3 {
		return firstOperand, secondOperand, op, errors.New("Missing parameter!!!, Example 2 + 3 ")
	}

	firstOperand, err = strconv.ParseFloat(inputParse[0], 10)
	if err != nil {
		return
	}
	if strings.HasSuffix(inputParse[2], "\n") {
		secondOperand, err = strconv.ParseFloat(strings.ReplaceAll(inputParse[2], "\n", ""), 10)
	} else {
		secondOperand, err = strconv.ParseFloat(inputParse[2], 10)
	}
	if err != nil {
		return
	}
	op = inputParse[1]

	return

}

func contains(arraySource []string, object string) bool {
	for _, v := range arraySource {
		if v == object {
			return true
		}
	}
	return false
}
