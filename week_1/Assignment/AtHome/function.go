package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func calculate(input string) (result float64, err error) {
	inputSplit, err := splitInput(input)
	if err != nil {
		return
	}
	//check if input < 3 element -> missing as least argument
	//or check if number argument is odd, the 2 element same type
	if len(inputSplit) < 3 || len(inputSplit)%2 == 0 {
		return 0, errors.New(fmt.Sprintf("Invalid input"))
	}
	currentOp := "+"
	var operand Operand
	var operation Operation
	for i, v := range inputSplit {
		//if i is odd the value is operand else the value is operation
		if i%2 == 0 {
			operand = Operand{orgInput: v}
			_, err := operand.check()
			if err != nil {
				return 0, err
			}
			result, err = operationHandle(result, operand.value, currentOp)
			if err != nil {
				return 0, err
			}
			continue
		} else {
			operation = Operation{orgInput: v}
			_, err := operation.check()
			if err != nil {
				return 0, err
			}
			currentOp = operation.op
			continue
		}
	}
	return
}

func splitInput(input string) (result []string, err error) {
	input, err = replaceAllSubString(input, " ", " ")
	if err != nil {
		return
	}
	return strings.Split(strings.Trim(input, " "), " "), nil
}

func replaceAllSubString(input string, subString string, newString string) (result string, err error) {
	regexSpace, err := regexp.Compile(`\s+`)
	if err != nil {
		return
	}
	return regexSpace.ReplaceAllString(input, newString), nil
}

func operationHandle(firstOperand, secondOperand float64, op string) (result float64, err error) {
	switch op {
	case "+":
		return firstOperand + secondOperand, nil
	case "-":
		return firstOperand - secondOperand, nil
	case "*":
		return firstOperand * secondOperand, nil
	case "/":
		if secondOperand == 0 {
			return 0, errors.New(fmt.Sprintf("Not divide by 0"))
		}
		return firstOperand / secondOperand, nil
	default:
		return 0, errors.New(fmt.Sprintf("Not support op: %s", op))
	}
}

type Checker interface {
	check() (value interface{}, err error)
}

type Operation struct {
	orgInput string
	op       string
}

type Operand struct {
	orgInput string
	value    float64
}

func (op *Operation) check() (value string, err error) {
	if contains(supportOperations, op.orgInput) {
		op.op = op.orgInput
		return op.op, nil
	}
	return "", errors.New(fmt.Sprintf("Not support op: %s", op.orgInput))
}

func (o *Operand) check() (value float64, err error) {
	value, err = strconv.ParseFloat(o.orgInput, 10)
	if err != nil {
		return
	}
	o.value = value
	return
}
