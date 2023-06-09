## Setting Defaults

You can set default struct values by unmarshalling the json bytes to that struct. Any new value from the json bytes will override the default values.
```golang
package main

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Name        string `json:"name"`
	Application string `json:"application"`
}

func main() {
	body := []byte(`{"name": "john"}`)

	def := Request{
		Name:        "test",
		Application: "web", // Setting this as default.
	}
	json.Unmarshal(body, &def)
	fmt.Println(def)
}
```

## Returning empty array instead of null in JavaScript

 ```golang
 var data []string // This will return null after unmarshalling
 data := make([]string, 0) // This will return empty array, []
 data = nil // Setting this back to nil will return null at the end, so if there are functions that returns nil, do a checking at the end
 
 if data == nil {
  data = make([]string, 0)
 }
 ```


## JSON Stream

When working with large json (e.g. few MBs to few GBs), using `json.NewDecoder` can be more performant than `json.Unmarshal` since we can partially unmarshal and batch validate the json before further processing.


```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

var raw = []byte(`{
	"data": [
		{"name": "john"},
		{"name": "jane"}
	],
	"meta": {}
}`)

type Person struct {
	Name string `json:"name"`
}

func main() {
	var result []Person
	dec := json.NewDecoder(bytes.NewReader(raw))
	mustDelim(dec, '{')
	for dec.More() {
		field, err := dec.Token()
		if err != nil {
			panic(err)
		}
		if field != "data" {
			continue
		}
		mustDelim(dec, '[')
		for dec.More() {
			var p Person
			if err := dec.Decode(&p); err != nil {
				panic(err)
			}
			result = append(result, p)
		}
		mustDelim(dec, ']')
	}
	mustDelim(dec, '}')
	fmt.Println("person", result)
}

func mustDelim(dec *json.Decoder, delim json.Delim) {
	token, err := dec.Token()
	if err != nil {
		panic(err)
	}
	del, ok := token.(json.Delim)
	if !ok || del != delim {
		log.Fatalf("expected %v, got %v", delim, del)
	}
}
```


## Extracting and validating if a json string is valid

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

var text = `hello

{
	"name": "john"
}
`

func main() {
	a := strings.Index(text, "{")
	b := strings.LastIndex(text, "}")
	res := text[a : b+1]
	fmt.Println(res)
	fmt.Println(json.Valid([]byte(res)))
}
```
