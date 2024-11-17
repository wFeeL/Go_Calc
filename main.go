package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	tokens, err := EvaluateSign(expression)
	if err != nil {
		return 0, err
	}
	return ParseValues(tokens)
}

func EvaluateSign(expression string) ([]string, error) {
	var tokens []string
	var number strings.Builder

	for _, ch := range expression {
		if unicode.IsSpace(ch) {
			continue
		}
		if unicode.IsDigit(ch) || ch == '.' {
			number.WriteRune(ch)
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			if strings.ContainsRune("+-*/()", ch) {
				tokens = append(tokens, string(ch))
			} else {
				return nil, errors.New("invalid character in expression")
			}
		}
	}
	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}
	return tokens, nil
}

func ParseValues(tokens []string) (float64, error) {
	var values []float64
	var operators []string

	var applyOperator = func() error {
		if len(values) < 2 || len(operators) == 0 {
			return errors.New("invalid expression")
		}
		right, left := values[len(values)-1], values[len(values)-2]
		values = values[:len(values)-2]

		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		var result float64
		switch op {
		case "+":
			result = left + right
		case "-":
			result = left - right
		case "*":
			result = left * right
		case "/":
			if right == 0 {
				return errors.New("division by zero")
			}
			result = left / right
		default:
			return errors.New("unknown operator")
		}
		values = append(values, result)
		return nil
	}

	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			values = append(values, num)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				if err := applyOperator(); err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 || operators[len(operators)-1] != "(" {
				return 0, errors.New("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]
		} else if strings.Contains("+-*/", token) {
			for len(operators) > 0 && GetPrioritiesOperation(operators[len(operators)-1]) >= GetPrioritiesOperation(token) {
				if err := applyOperator(); err != nil {
					return 0, err
				}
			}
			operators = append(operators, token)
		} else {
			return 0, errors.New("unknown token")
		}
	}

	for len(operators) > 0 {
		if err := applyOperator(); err != nil {
			return 0, err
		}
	}

	if len(values) != 1 {
		return 0, errors.New("invalid expression")
	}
	return values[0], nil
}

func GetPrioritiesOperation(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}
