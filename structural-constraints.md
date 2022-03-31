# Structural constraints in go 1.18

We can now enforce struct fields, but for different types. The disadvantage is all the fields must be exact:

```go

// You can edit this code!
// Click here and start typing.
package main

import "fmt"

type User struct {
	ID string
}

type UserAPI struct {
	ID string
}

func main() {
	Print(User{ID: "user"})
	Print(UserAPI{ID: "user-api"})
}

type HasID interface {
	~struct {
		ID string
	}
}

func Print[T HasID](t T) {
	fmt.Println(t.ID)
}
```
