## Specification pattern with golang


```go
package main

import (
	"fmt"
)

type Specification interface {
	IsSatisfiedBy(User) bool

	// Optionally...
	// And(Specification) bool
	// Or(Specification) bool
	// Not(Specification) bool
}

type User struct {
	Age    int
	Gender string
}

type AgeSpecification struct {
	threshold int
}

func (a *AgeSpecification) IsSatisfiedBy(user User) bool {
	return user.Age >= a.threshold
}

func NewAgeSpecification(threshold int) *AgeSpecification {
	return &AgeSpecification{threshold: threshold}
}

type GenderSpecification struct {
	gender string
}

func (g *GenderSpecification) IsSatisfiedBy(user User) bool {
	return g.gender == user.Gender
}

func NewGenderSpecification(gender string) *GenderSpecification {
	return &GenderSpecification{gender: gender}
}

func main() {
	ageSpec := NewAgeSpecification(13)
	genderSpec := NewGenderSpecification("male")

	user := User{Age: 20, Gender: "male"}

	if valid := ageSpec.IsSatisfiedBy(user); valid {
		fmt.Println(valid)
	}

	if valid := genderSpec.IsSatisfiedBy(user); valid {
		fmt.Println(valid)
	}

}
```
