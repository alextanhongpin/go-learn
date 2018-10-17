Write a program to print students in format 'name-sex'. The conditions are:
- boys first, followed by girl
- you can only use on for/while loop (either one), and nested for/while loop is not allowed
- only basic data type can be used such as: int, float, double, bool
- array, list, string, and stringbuffer is not allowed

```go
package main

import (
	"fmt"
)

type Student struct {
	name string
	sex  rune
}

func main() {
	students := []Student{
		{"Martin", 'M'},
		{"Jenny", 'F'},
		{"Wendy", 'F'},
		{"Siti", 'F'},
		{"Thomas", 'M'},
		{"Siva", 'M'},
		{"Mun Hui", 'F'},
		{"Richard", 'M'},
		{"Kumar", 'M'},
		{"Isah", 'F'},
		{"Samson", 'M'},
	}
	n := len(students)

	gender := 'M'
	for i := 0; i < n*2; i++ {
		if i >= n {
			gender = 'F'
		}
		if student := students[i%n]; student.sex == gender {
			fmt.Printf("%s-%c\n", student.name, student.sex)
		}
	}
}
```

Alternative solution is to use queue, pop the first item, print it if it is male. Else push it back at the end of the array. Then print the rest as female.
