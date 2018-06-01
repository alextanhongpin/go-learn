package main

import (
	"fmt"
)

var days = [...]string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}

type Day int

const (
	Sunday Day = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (d Day) String() string {
	return days[d]
}

func main() {
	fmt.Printf("Today is %s.\n", Sunday)
	fmt.Printf("Type Sunday: %#v\n", Sunday)

}
