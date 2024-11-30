
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

## With MarshalJSON
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type DateJSON struct {
	time.Time
}

func (d DateJSON) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return d.Time.MarshalJSON()
}

type BirthDate struct {
	time.Time
}

func (b BirthDate) MarshalJSON() ([]byte, error) {
	if b.IsZero() {
		return []byte("null"), nil
	}
	return b.Time.MarshalJSON()
}

func (b BirthDate) Age(now time.Time) int {
	year, month, date, _, _, _ := diff(b.Time, now)
	// More than half a month is considered a month.
	if date > 15 {
		month++
	}
	// More than half a year is one year.
	if month > 6 {
		year++
	}
	return year
}

type User struct {
	Name      string    `json:"name"`
	CreatedAt DateJSON  `json:"created_at,omitempty"`
	BirthDate BirthDate `json:"birth_date,omitempty"`
}

func main() {
	// From JSON.
	{
		js := `{"name":"john","created_at":"2009-11-10T23:00:00Z", "birth_date": "1990-01-01T00:00:00Z"}`
		var u User
		err := json.Unmarshal([]byte(js), &u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(u.CreatedAt)
		fmt.Println(u.BirthDate, u.BirthDate.Age(time.Date(2019, 5, 13, 0, 0, 0, 0, time.Local)))
	}
	// To JSON
	{
		user := &User{
			Name:      "john",
			CreatedAt: DateJSON{time.Now()},
		}
		b, _ := json.Marshal(user)
		fmt.Println(string(b))
	}
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
```


## Skipping business days

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Date(2021, 6, 20, 0, 0, 0, 0, time.Local)
	fmt.Println("today is", start.Weekday())
	delivered := estimateDelivery(start, 3)
	fmt.Println("expected to receive on ", delivered.Weekday())
}

func estimateDelivery(start time.Time, days int) time.Time {
	for days > 0 {
		start = start.AddDate(0, 0, 1)
		day := start.Weekday()
		if day == time.Saturday || day == time.Sunday {
			start = start.AddDate(0, 0, 1)
		} else {
			days--
		}
	}
	return start
}
```

Without loop:
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().AddDate(0, 0, 3)
	fmt.Println("Start", now.Weekday())
	end := AddBusinessDays(now, 6)
	fmt.Println(end.Weekday(), end)
}

// https://stackoverflow.com/questions/1044688/addbusinessdays-and-getbusinessdays
func AddBusinessDays(start time.Time, days int) time.Time {
	if start.Weekday() == time.Saturday {
		start = start.AddDate(0, 0, 2)
		days -= 1
	} else if start.Weekday() == time.Sunday {
		start = start.AddDate(0, 0, 1)
		days -= 1
	}
	initialDayOfWeek := int(start.Weekday())
	weeks := days / 5
	addDays := days % 5
	if addDays+initialDayOfWeek > 5 {
		addDays += 2
	}
	return start.AddDate(0, 0, weeks*7+addDays)
}
```


## Start and End of day

```go
now := time.Now().In(model.GMTe7)
startOfDay := now.Add(-12 * time.Hour).Round(24 * time.Hour)
endOfDay := now.Add(24 * time.Hour)
```

## Exponential Decay for effective polling


```go

// combination of two curves. the duration increases exponentially in the beginning before beginning to decay.
// The idea is the wait duration should eventually be lesser and lesser over time.
func exponentialGrowthDecay(i int) time.Duration {
	x := float64(i)
	base := 1.0 + rand.Float64()
	switch {
	case x < 4: // intersection point rounded to 4
		base *= math.Pow(2, x)
	case x < 10:
		base *= 5 * math.Log(-0.9*x+10)
	default:
	}

	return time.Duration(base*100) * time.Millisecond
}
```
