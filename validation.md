## Validation strategy

- validate externally (for simple logic, primary types etc)
- using method structs (used to wrap common validation logic such as required fields etc)
- wrap another struct with a validator (when the struct to validate has different validation logic for different scenarios and they don't share the same validation logic, normally this happens for request/response Data Transfer Object, dto, e.g. admin requires a specific field etc)
- use third party validator with tags (most flexible, but sometimes we want to avoid third party library for simple validation)

## Validate externally

```go

i := 100
if i > 100 {
  // do something
}
```

## Using struct methods

```go
package main

import (
	"errors"
	"log"
)

type ListingRequest struct {
	Limit int
}

func (l *ListingRequest) Validate() error {
	if l.Limit > 100 {
		return errors.New("limited to 100 per page")
	}
	return nil
}

func main() {
	req := &ListingRequest{999}
	if err := req.Validate(); err != nil {
		log.Fatal(err)
	}
}
```

## Wrap another struct with a validator

```go
package main

import (
	"errors"
	"log"
)

type ListingRequest struct {
	Limit int
}

type AuthenticatedListingRequestValidator struct {
	request *ListingRequest
}

func (a AuthenticatedListingRequestValidator) Validate() error {
	return nil
}

type PublicListingRequestValidator struct {
	request *ListingRequest
}

func (p *PublicListingRequestValidator) Validate() error {
	if p.request.Limit > 1000 {
		return errors.New("limited to 100 per page")
	}
	return nil
}

func main() {
	req := &ListingRequest{2000}
	err := (&AuthenticatedListingRequestValidator{req}).Validate()
	if err != nil {
		log.Fatal(err)
	}
	err = (&PublicListingRequestValidator{req}).Validate()
	if err != nil {
		log.Fatal(err)
	}
}
```

Or:

```go
package main

import (
	"fmt"
	"log"
)

type ListingRequest struct {
	Limit int
}

type ListingRequestValidator struct {
	threshold int
}

func (l *ListingRequestValidator) Validate(req *ListingRequest) error {
	if req.Limit > l.threshold {
		return fmt.Errorf(`limited to %d per page`, l.threshold)
	}
	return nil
}

func main() {
	req := &ListingRequest{999}
	authValidator := &ListingRequestValidator{1000}
	publicValidator := &ListingRequestValidator{100}
	err := authValidator.Validate(req)
	if err != nil {
		log.Fatal(err)
	}
	err = publicValidator.Validate(req)
	if err != nil {
		log.Fatal(err)
	}
}
```


## Validation 

How can we keep our data separate from behaviours? The advantage of programming is it is flexible. However, sometimes we just want to ensure that validation could be represented as a data instead of code.

Also for errors, most of the time, we want to
1) identify exactly what happens through a meaningful error code
2) pass additional data down to construct more meaningful error messages (as well as localization)
3) validation may depend on external dependencies (another entity, some external infrastructure, service or third-party APIs)

```go
package main

import (
	"errors"
	"fmt"
)

var ErrNameTooLong = errors.New("user.nameTooLong")
var ErrNameIsRequired = errors.New("user.nameIsRequired")

type User struct {
	Name string
	Age  int
}

// Using interface without arguments in function keeps the method as generic as possible.
type Validator interface {
	Validate() error
}

type UserNameValidator struct {
	user *User
}

func (v *UserNameValidator) Validate() error {
	isEmpty := v.user.Name == ""
	isTooLong := len(v.user.Name) > 100

	switch {
	case isEmpty:
		return ErrNameIsRequired
	case isTooLong:
		return ErrNameTooLong
	default:
		return nil
	}
}

func main() {
	u := &User{Name: "John"}
	validators := []Validator{
		&UserNameValidator{u},
		// TODO: Add validator that can backfill empty values.
	}
	for _, v := range validators {
		if err := v.Validate(); err != nil {
			panic(err)
		}
	}
	fmt.Println("Hello, playground")
}

```
