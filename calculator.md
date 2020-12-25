## Calculator

Basic calculator implementation using [operatator precedence parser](https://en.wikipedia.org/wiki/Operator-precedence_parser).

It only works when:
- each numeric token is between 0-9
- no brackets

```go
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var operators = map[string]int{
	"-": 1,
	"+": 1,
	"*": 2,
	"/": 2,
}

func main() {
	input := "2+9*2/6"
	tokens := strings.Split(input, "")

	var i int
	advance := func() {
		i++
	}

	parsePrimary := func() (int, bool) {
		if i >= len(tokens) {
			return 0, false
		}
		n, err := strconv.Atoi(tokens[i])
		if err != nil {
			log.Fatal(err)
		}
		return n, true
	}

	peek := func() string {
		if i+1 >= len(tokens) {
			return ""
		}
		return tokens[i+1]
	}

	precedence := func(s string) int {
		return operators[s]
	}

	operation := func(op string, lhs, rhs int) int {
		switch op {
		case "+":
			return lhs + rhs
		case "-":
			return lhs - rhs
		case "*":
			return lhs * rhs
		case "/":
			return lhs / rhs
		}
		return 0
	}

	var eval func(lhs, minPrecedence int) int
	eval = func(lhs, minPrecedence int) int {
		lookahead := peek()
		for precedence(lookahead) >= minPrecedence {
			advance()
			op := lookahead
			
			advance()
			rhs, ok := parsePrimary()
			if !ok {
				break
			}
			lookahead = peek()
			for precedence(lookahead) > precedence(op) {
				rhs = eval(rhs, precedence(lookahead))
				lookahead = peek()
			}
			lhs = operation(op, lhs, rhs)
		}
		return lhs
	}
	lhs, _ := parsePrimary()
	fmt.Println(eval(lhs, 0))
}
```


## Calculator v2

Using two stacks to handle operations with brackets, read [here](http://waltermilner.com/opp.pdf).

Two stacks:
- numbers: stores the numeric value
- operators: stores the operators, including brackets


Approach:
1. Parse the token from left to right
2. If it is a number, push to the `numbers` stack
3. If it is a left bracket, push to the `operators` stack
4. If it is a right bracket `)`
	- `reduce` until the top operator is a left bracket `(`
	- pop the left bracket
5. Else it must be an operator, `currOp`
	- peek at the top of the `operators` stack, `lastOp`
	- if `lastOp` has higher precedence than `currOp` (e,g. `*` multiplication has higher precedence than `+` addition), then do a `reduce`
	- push the `currOp` into the `operators` stack

A `reduce` is merely popping out two values from the `numbers` (`rhs`, followed by `lhs`) stack, and one value from the `operators` stack, and performing the operation `lhs op rhs`.


```go
package main

import (
	"log"
	"strconv"
	"strings"
)

var operators = map[string]int{
	"-": 1,
	"+": 1,
	"*": 2,
	"/": 2,
}

func main() {
	assert := func(in string, out int) {
		if n := eval(in); n != out {
			log.Fatalf("expected %s=%d, got %d", in, out, n)
		}
	}
	assert("1", 1)
	assert("1+1", 2)
	assert("1+2-3", 0)
	assert("1+2*3", 7)
	assert("1+2*3-(4-2)", 5)
	assert("5*2-(5-2)", 7)
}

func eval(in string) int {
	tokens := strings.Split(in, "")

	var ops []string
	var numbers []int
	var token string

	reduce := func() {
		var lhs, rhs int
		var op string
		op, ops = ops[len(ops)-1], ops[:len(ops)-1]
		rhs, numbers = numbers[len(numbers)-1], numbers[:len(numbers)-1]
		lhs, numbers = numbers[len(numbers)-1], numbers[:len(numbers)-1]

		numbers = append(numbers, operation(op, lhs, rhs))
	}
	
	for len(tokens) > 0 {
		token, tokens = tokens[0], tokens[1:]
		if token == "(" {
			ops = append(ops, token)
		} else if token == ")" {
			for ops[len(ops)-1] != "(" {
				reduce()
				if ops[len(ops)-1] == "(" {
					ops = ops[:len(ops)-1]
					break
				}
			}
		} else if checkOperator(token) {
			t, ok := tail(ops)
			if ok && precedence(t) > precedence(token) {
				reduce()
			}
			ops = append(ops, token)
		} else {
			numbers = append(numbers, toInt(token))
		}
	}
	for len(ops) > 0 {
		reduce()
	}
	return numbers[0]
}

func precedence(s string) int {
	return operators[s]
}

func operation(op string, lhs, rhs int) int {
	switch op {
	case "+":
		return lhs + rhs
	case "-":
		return lhs - rhs
	case "*":
		return lhs * rhs
	case "/":
		return lhs / rhs
	}
	return 0
}

func checkOperator(s string) bool {
	_, exists := operators[s]
	return exists
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func head(arr []string) (string, bool) {
	if len(arr) == 0 {
		return "", false
	}
	return arr[0], true
}

func tail(arr []string) (string, bool) {
	if len(arr) == 0 {
		return "", false
	}
	return arr[len(arr)-1], true
}
```
