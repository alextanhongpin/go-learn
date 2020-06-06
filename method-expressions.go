```go
package main

import (
	"fmt"
)

type User struct {
}

func (u User) ValueReceiver(name string) {
	fmt.Println(name)
}
func (u *User) PointerReceiver(name string) {
	fmt.Println(name)
}

func main() {
	// Method expressions.
	var u User
	u.ValueReceiver("u.ValueReceiver")
	u.PointerReceiver("u.PointerReceiver")
	(User).ValueReceiver(u, "(User).ValueReceiver")
	User.ValueReceiver(u, "(User).ValueReceiver")
	(*User).PointerReceiver(&u, "(User).PointerReceiver")
}
```
