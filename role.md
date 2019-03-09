## Role-based

When implementing RBAC or any role-based permission for APIs, we sometimes need to hide certain fields depending on the owner of the system. We can use type conversions to hide the json tags that are returned to the client.

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Base struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Passport string `json:"passport"`
}

func (b Base) Info() {
	fmt.Printf("%s (age %d) passport is %s\n", b.Name, b.Age, b.Passport)
}

type Private struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Passport string `json:"-"`
}

func (p Private) Secret() {
	fmt.Printf("%s passport's is kept a secret\n", p.Name)
}

func main() {
	user := Base{
		Name:     "john",
		Age:      10,
		Passport: "ABC",
	}
  // Only Owner can view this information.
	user.Info()
  
  // Other users can only view the limited response.
	private := Private(user)
	private.Secret()
	b, _ := json.Marshal(private)
	fmt.Println(string(b))
}
```
