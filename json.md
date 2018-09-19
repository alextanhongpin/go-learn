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
