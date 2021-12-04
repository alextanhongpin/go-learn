
## Complex Number

We can use complex number to represent position and direction. Useful for puzzle games where you need to navigate up/down/left/right and switch direction.

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {

	// The real number is the x-coordinate.
	// The imaginary number is the y-coordinate.
	pos := 0 + 0i // (x: 0, y: 0)
	fmt.Println(pos)

	// Adding real number moves the position to the right, and vice-versa.
	fmt.Println(pos+1, pos-1)

	// Adding imaginary number moves the position to the bottom, and vice-versa.
	fmt.Println(pos+1i, pos-1i)

	dir := 0 + 1i
	fmt.Println("top:", dir)            // Facing top
	fmt.Println("left:", dir*dir)       // Facing left (rotate -90deg)
	fmt.Println("bottom:", dir*dir*dir) // Facing bottom
	fmt.Println("bottom:", dir*-1)
	fmt.Println("right:", dir*dir*dir*dir)   // Facing right
	fmt.Println("right:", dir*dir*-1)        // Also facing right.
	fmt.Println("top:", dir*dir*dir*dir*dir) // Facing top

	// Move top
	fmt.Println(pos+dir, pos+dir+dir)

	// Move left.
	fmt.Println(pos + (dir * dir))
}
```
