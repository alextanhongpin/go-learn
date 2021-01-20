# Stringer
```go
package main

import "fmt"

type PokemonType int

const (
	PokemonTypeFire PokemonType = iota
	PokemonTypeWater
	PokemonTypeGrass
)

type Day int

const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type Card int

const (
	CardAce Card = iota
	CardKing
)

type Name int

const (
	James Name = iota // james
	Alice             // alice
)

// go get https://godoc.org/golang.org/x/tools/cmd/stringer
//go:generate stringer -type=PokemonType
//go:generate stringer -type=Day
//go:generate stringer -type=Card -trimprefix=Card
//go:generate stringer -type=Name -linecomment
func main() {
	fmt.Println(PokemonTypeFire)
	fmt.Println(Sunday)
	fmt.Println(CardAce)
	fmt.Println(James)
}
```
