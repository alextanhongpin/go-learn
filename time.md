
## Zero time

This program demonstrates how to handle zero time. When returning the time as JSON, use pointer to avoid zero time in JSON response.

```go
package main

import (
	"log"
	"time"
)

type Book struct {
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func main() {

	var b Book
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = now
	}

	log.Println(b)
}
```

## Local time
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	singapore := time.FixedZone("Singapore Time", secondsEastOfUTC)
	fmt.Println(time.Now(), time.Now().In(singapore).Format(time.RFC3339))
}
```


## Timezone

```go
package main

import (
	"log"
	"strings"
	"time"
)

func main() {
	{
		t := time.Now()
		zone, offset := t.Zone()
		log.Println(zone, offset)
	}
	{
		t := time.Now().In(time.Local)
		zone, offset := t.Zone()
		log.Printf("%+v, %+v\n", zone, offset)
		log.Println(t.Location().String())
	}
	{
		t, _ := time.Parse(time.RFC3339, "2019-04-09T11:34:20+08:00")
		zone, offset := t.Zone()
		log.Println(zone, offset)
		log.Println(t.Location().String())
	}
	{
		dateArray := strings.Fields(time.Now().String())
		log.Println("date", dateArray[0])
		log.Println("time", dateArray[1])
		log.Println("offset", dateArray[2])
		log.Println("timezone", dateArray[3])
	}
}
```

## Start and end of the month

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	firstday := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	lastday := firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
	fmt.Println(t)
	fmt.Println(firstday)
	fmt.Println(lastday)
}
```

## Start of Day

```go
func startOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
```

## Removing JSON Date

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type DateOfBirth time.Time

func NewDateOfBirth(year, month, day int) DateOfBirth {
	return DateOfBirth(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local))
}

func (d DateOfBirth) IsValid() bool {
	return !time.Time(d).Equal(time.Date(9999, 12, 31, 0, 0, 0, 0, time.Local)) && !time.Time(d).IsZero()
}
func (d DateOfBirth) Date() time.Time {
	return time.Time(d)
}

type User struct {
	Name            string      `json:"name"`
	DateOfBirth     DateOfBirth `json:"-"`
	DateOfBirthJSON *time.Time  `json:"date_of_birth,omitempty"`
}

func main() {
	u := &User{
		Name: "john",
		// DateOfBirth: NewDateOfBirth(1990, 1, 1),
	}
	if u.DateOfBirth.IsValid() {
		var date = u.DateOfBirth.Date()
		fmt.Println(date)
		u.DateOfBirthJSON = &date
	}
	b, _ := json.Marshal(u)
	fmt.Println(string(b))
}
```
