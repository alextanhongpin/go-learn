Sometimes we need to deal with slightly more dynamic data structure (dynamic json forms etc), but we still want to take advantage of the strongly typed feature of golang, and the struct methods. We can implement a mapper that converts the concrete struct into dynamic values and vice-versa.
```golang
package main

import (
	"fmt"
)

type Generic struct {
	Key   string      // The field name
	Value interface{} // The value, can be cast back to the original data type.
	Type  string      // The data type
}

// Concrete type to generic type conversion
type Concrete struct {
	Subject string
	Body    string
}

func (c *Concrete) Values() []Generic {
	return []Generic{
		{"subject", c.Subject, "string"},
		{"body", c.Body, "string"},
	}
}

func (c *Concrete) From(generics []Generic) error {
	for _, gen := range generics {
		switch gen.Key {
		case "subject":
			subject, _ := gen.Value.(string)
			c.Subject = subject
		case "body":
			body, _ := gen.Value.(string)
			c.Body = body
		default:
			return fmt.Errorf("error parsing field: %s", gen.Key)
		}
	}
	return nil
}

func main() {
	c := Concrete{
		Subject: "a",
		Body:    "b",
	}
	fmt.Println(c.Values())
	var d Concrete
	d.From(c.Values())
	fmt.Println(d)
}

```
