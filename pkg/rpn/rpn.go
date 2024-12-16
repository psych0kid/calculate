package rpn

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Calc evaluates a mathematical expression given in Reverse Polish Notation (RPN)
func Calc(expression string) (float64, error) {
	prepared := preparingExpression(expression)
	rpn, err := convertToRPN(prepared)
	if err != nil {
		return 0, err
	}
	return calculateRPN(rpn)
}

// preparingExpression adds a zero before a negative sign if it is the first character of the expression.
func preparingExpression(expression string) string {
	var prepared strings.Builder
	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])
		if char == '-' {
			if i == 0 {
				prepared.WriteRune('0')
			} else if rune(expression[i-1]) == '(' {
				prepared.WriteRune('0')
			}
		}
		prepared.WriteRune(char)
	}
	return prepared.String()
}

// convertToRPN converts an infix expression to Reverse Polish Notation (RPN).
func convertToRPN(expression string) ([]string, error) {
	var output []string
	var operators []rune
	var priority = map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	unary := true
	sizeOpr := -1

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if unicode.IsDigit(char) || char == '.' || (char == '-' && unary) {
			var number strings.Builder
			if char == '-' && unary {
				number.WriteRune(char)
				i++
				char = rune(expression[i])
			}
			number.WriteRune(char)
			for i+1 < len(expression) && (unicode.IsDigit(rune(expression[i+1])) || rune(expression[i+1]) == '.') {
				i++
				number.WriteRune(rune(expression[i]))
			}

			if sizeOpr != len(operators) {
				output = append(output, number.String())
				unary = false
				sizeOpr = len(operators)
			} else {
				return nil, errors.New("error expression")
			}
		} else if char == '(' {
			operators = append(operators, char)
			unary = true
		} else if char == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 || operators[len(operators)-1] != '(' {
				return nil, errors.New("mismatched brackets")
			}
			operators = operators[:len(operators)-1]
			unary = false
		} else if strings.ContainsRune("-+*/", char) {
			for len(operators) > 0 && priority[operators[len(operators)-1]] >= priority[char] {
				output = append(output, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)
			unary = true
			sizeOpr = len(operators) - 1
		} else if !unicode.IsSpace(char) {
			return nil, errors.New("invalid character in expression")
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == '(' {
			return nil, errors.New("mismatched brackets")
		}
		output = append(output, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

// calculateRPN evaluates a Reverse Polish Notation (RPN) expression.
func calculateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else if len(stack) >= 2 {
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "-":
				stack = append(stack, a-b)
			case "+":
				stack = append(stack, a+b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("division on zero")
				}
				stack = append(stack, a/b)
			default:
				return 0, errors.New("unknown operator")
			}
		} else {
			return 0, errors.New("error in expression")
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("error in expression")
	}

	if math.Abs(stack[0]) < 1e-9 {
		return 0, nil
	}

	return stack[0], nil
}
