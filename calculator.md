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
