```go
package main

import (
	"fmt"
	"log"
	"time"
)

type DateParser struct {
	layout string
}

func (d *DateParser) Parse(str string) (time.Time, error) {
	return time.Parse(d.layout, str)
}

func (d *DateParser) Format(t time.Time) string {
	return t.Format(d.layout)
}
func NewISO8601DateParser() *DateParser {
	return &DateParser{"2006-01-02"}
}


func main() {
	parser := NewISO8601DateParser()
	t, err := parser.Parse("1990-01-01")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)
	fmt.Println(parser.Format(time.Now()))
}
```


```go
package main

import (
	"fmt"
	"time"
)

type MonthRange struct {
	time time.Time
}

func (m MonthRange) Start() time.Time {
	year, month, _ := m.time.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, m.time.Location())
}

func (m MonthRange) End() time.Time {
	return m.Start().AddDate(0, 1, 0).Add(time.Nanosecond * -1)
}

func main() {
	mth := MonthRange{time: time.Now()}
	fmt.Println(mth.Start(), mth.End())
}
```
