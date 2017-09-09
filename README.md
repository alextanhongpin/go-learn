# 10 Things I learnt about Go

## 1. There's no round function for numbers

If you want to round numbers in go, you have to implement it yourself, that is, until [go 1.10](https://github.com/golang/go/issues/20100) is released.

```golang
// This program demonstrates how to round numbers in golang

package main

import (
	"log"
	"math"
)

func round(val float64) float64 {
	_, v := math.Modf(val)
	if math.Abs(v) >= .5 {
		if val >= 0 {
			return math.Ceil(val)
		}
		return math.Floor(val)
	}
	if val >= 0 {
		return math.Floor(val)
	}
	return math.Ceil(val)
}

func main() {
	log.Println(round(123.54))
	log.Println(round(-100.4))
	log.Println(round(-0.4))
	log.Println(round(-0.5))
	log.Println(round(0.5))
	log.Println(round(0.4))
}
```

## 2. There's no reverse string function

Again, it's weird that `go` lacks something as simple as a reverse string function. But there are times when you actually need it. Here's how you can do it:

```golang
// This program demonstrates how to reverse string in go
package main

import "log"

func main() {
	log.Println(reverse("hello world!"))
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

```

## 3. Writing to a JSON file

Marshalling is the process of converting domain objects to a serialized format, such as `json`. In order to write a `golang` struct to a `json` file, you need to marshal it first:

```golang
// This program demonstrates how to write a struct to json file

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Point represents the schema of our json output
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func writeJSON(file string, obj interface{}, pretty bool) (err error) {
	var bytes []byte
	if pretty {
		bytes, err = json.MarshalIndent(obj, "", "  ")
	} else {
		bytes, err = json.Marshal(obj)
	}
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, bytes, 0644)
}

func main() {
	points := []Point{Point{0, 0}, Point{1, 1}}
	err := writeJSON("out.json", points, true)
	if err != nil {
		log.Fatalf("error writing to json: %v\n", err)
	}
}
```

If you set `pretty` to `true`, it will be saved in a more readable format. This is how our output `json` will look like:
```json
[
  {
    "x": 0,
    "y": 0
  },
  {
    "x": 1,
    "y": 1
  }
]
```

## 4. Loading a JSON file

Unmarshalling is the process of converting domain objects from a serialized format, such as `json`. To load the json file to our `golang` struct, we have to unmarshal the `json` data:

```golang
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func loadJSON(file string, obj interface{}) error {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &obj)
}

// Point represents the schema of the json we want to load
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	var points []Point
	if err := loadJSON("out.json", &points); err != nil {
		log.Printf("error loading json: %v", err)
	}
	log.Printf("load json: %#v\n", points)
}
```

This is the `json` file we are loading:

```json
[
  {
    "x": 0,
    "y": 0
  },
  {
    "x": 1,
    "y": 1
  }
]
```


## 5. Mapping map to structs

In case you need to map golang `map` to `structs`, there is a library for it:

```golang
package main

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

type People struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var people map[string]interface{}

func main() {
	// Using Hashicorp's library to convert map to struct
	// Example struct
	people = make(map[string]interface{})
	people["name"] = "car" // Lowercase works
	people["Age"] = 1

	peeps := People{}
	err := mapstructure.Decode(people, &peeps)
	if err != nil {
		log.Println(err)
	}
	log.Printf("peeps: %#v\n", peeps)
}
```

## 6. Shadowing fields

There are times where you want to hide certain fields from the golang `struct` before returning it as a `json` response, but not with the `json:"-"` approach. The example below shows how you remove the *password* field from the original struct:

```golang
// This program demonstrates how to exclude fields in the json output
package main

import (
	"encoding/json"
	"log"
)

type UserPrivate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPublic struct {
	*UserPrivate
	Password bool `json:"password,omitempty"`
}

func main() {

	usrPriv := UserPrivate{"john.doe@mail.com", "123456"}
	usrPub := UserPublic{
		UserPrivate: &usrPriv,
	}
	// Convert it to bytes
	out, err := json.Marshal(usrPub)
	if err != nil {
		log.Println(err)
	}
	log.Printf("with shadowing: %s\n", string(out))
}
```

## 7. Composing struct

When returning a `json` response, you might want to return fields from other `structs`, but want to avoid creating too many of them. One way to achieve this is by composition - you compose a new struct by embedding other `structs` and you choose to exclude the fields too through __shadowing__ (see previous example).

```golang
// This program demonstrates how to compose multiple struct to be returned in the json
package main

import (
	"encoding/json"
	"log"
)

type User struct {
	Email    string `json:"email"`
}

type Skill struct {
	Name string `json:"name"`
	Level int `json:"level"`
}

type Skills []Skill

func main() {

	usr := User{"john.doe@mail.com"}
	skills := Skills{Skill{"javascript", 1}, {"go", 2}}

	// Convert our composed anyonymous struct to bytes
	out, err := json.Marshal(struct {
		*User
		*Skills `json:"skills"`
	} {
		User: &usr,
		Skills: &skills,
	})
	
	if err != nil {
		log.Println(err)
	}
    // Our json will have the skills field embedded
	log.Printf("with shadowing: %s\n", string(out))
}
```

## 8. Overwriting tag names

The name of the fields returned in the `json` is based on the json tag in your struct. You can overwrite them if you want your `json` response to have different field name:

```golang
// This program demonstrates how to overwrite the json tag in the struct
package main

import (
	"encoding/json"
	"log"
)

type BadField struct {
	Name string `json:"NameString"`
}

type GoodField struct {
	*BadField
	BadName string `json:"NameString,omitempty"`
	Name    string `json:"name"`
}

func main() {
	b := BadField{"john.doe"}
	g := GoodField{
		BadField: &b,
		Name:     b.Name,
	}

	out, err := json.Marshal(g)
	if err != nil {
		log.Printf("error unmarshalling: %v\n", err)
	}
	log.Println(string(out))
}
```

## 9. Concatenating arrays

It's probably wasn't that obvious, but concatenating array can be easily done as shown below:

```golang
// This program demonstrates how to concat two arrays in go

package main

import "log"

func main() {

	a := []string{"a", "b"}
	b := []string{"1", "2"}

	log.Println(append(a, b...))
}
```

Note that both arrays must be of the same type. Appending a `string` array to an `int` array will result in an error.

## 10. It's fast

Here's a benchmark of a "hello world" request using [wrk](https://github.com/wg/wrk). 

1 threads and 1 connections: 

```bash
# nodejs
wrk -d30s -c1 -t1 http://localhost:3000
Running 30s test @ http://localhost:3000
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    53.87us  201.49us  15.22ms   99.75%
    Req/Sec    19.70k     1.49k   20.53k    93.36%
  590138 requests in 30.10s, 62.47MB read
Requests/sec:  19606.26
Transfer/sec:      2.08MB

# go with standard library
wrk -d30s -c1 -t1 http://localhost:8080
Running 30s test @ http://localhost:8080
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    49.98us   30.88us   3.07ms   98.39%
    Req/Sec    19.26k     1.50k   20.21k    90.03%
  576548 requests in 30.10s, 70.38MB read
Requests/sec:  19154.60
Transfer/sec:      2.34MB

# go with fasthttp
wrk -d30s -c1 -t1 http://localhost:8080
Running 30s test @ http://localhost:8080
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    39.84us  138.55us   9.55ms   99.59%
    Req/Sec    27.32k     2.74k   32.76k    74.75%
  818105 requests in 30.10s, 113.91MB read
Requests/sec:  27179.54
Transfer/sec:      3.78MB
```

Similar test carried out with 10 threads and 10 connections: 

```bash
# nodejs
wrk -d30s -c10 -t10 http://localhost:3000
Running 30s test @ http://localhost:3000
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   360.79us  115.30us   4.22ms   87.49%
    Req/Sec     2.78k   586.91     3.18k    80.00%
  831578 requests in 30.10s, 88.03MB read
Requests/sec:  27627.49
Transfer/sec:      2.92MB

# go with standard library
wrk -d30s -c10 -t10 http://localhost:8080
Running 30s test @ http://localhost:8080
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   452.68us    2.30ms  96.75ms   97.89%
    Req/Sec     5.24k     1.06k    8.30k    75.03%
  1563840 requests in 30.02s, 190.90MB read
Requests/sec:  52090.76
Transfer/sec:      6.36MB

# go with fasthttp
wrk -d30s -c10 -t10 http://localhost:8080
Running 30s test @ http://localhost:8080
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   177.69us  510.81us  22.96ms   97.72%
    Req/Sec     7.82k     1.19k   17.26k    84.15%
  2336866 requests in 30.10s, 325.38MB read
Requests/sec:  77635.80
Transfer/sec:     10.81MB
```