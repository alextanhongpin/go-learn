## Composite

```go
package main

import (
	"fmt"
)

type Athlete struct{}

func (a *Athlete) Train() {
	fmt.Println("training")
}

type CompositeSwimmerA struct {
	MyAthlete *Athlete
	MySwim    func()
}

func Swim() {
	fmt.Println("swimming")
}

func main() {
	swimmer := CompositeSwimmerA{
		MySwim: Swim,
	}
	swimmer.MyAthlete.Train()
	swimmer.MySwim()

	trainAndSwimStruct(swimmer)
	trainAndSwim(swimmer.MyAthlete)
}

type Trainable interface {
	Train()
}

func trainAndSwimStruct(swimmer struct {
	MyAthlete *Athlete
	MySwim    func()
}) {
	swimmer.MyAthlete.Train()
	swimmer.MySwim()
}

func trainAndSwim(swimmer interface{ Trainable }) {
	swimmer.Train()
}
```

## Adapter 

```go
package main

import (
	"fmt"
)

type LegacyPrinter interface {
	Print(s string) string
}

type MyLegacyPrinter struct{}

func (l *MyLegacyPrinter) Print(s string) (newMsg string) {
	newMsg = fmt.Sprintf("legacy printer: %s\n", s)
	return
}

type ModernPrinter interface {
	PrintStored() string
}

type PrinterAdapter struct {
	OldPrinter LegacyPrinter
	Msg        string
}

func (p *PrinterAdapter) PrintStored() (newMsg string) {
	if p.OldPrinter != nil {
		newMsg = fmt.Sprintf("Adapter: %s", p.Msg)
		newMsg = p.OldPrinter.Print(newMsg)
		return
	}
	newMsg = p.Msg
	return
}

func main() {
	msg := "hello world"
	adapter := PrinterAdapter{
		OldPrinter: &MyLegacyPrinter{},
		Msg:        msg,
	}
	returnedMsg := adapter.PrintStored()
	fmt.Println(returnedMsg)

	newAdapter := PrinterAdapter{
		Msg: msg,
	}
	returnedMsg = newAdapter.PrintStored()
	fmt.Println(returnedMsg)
}
```

## Bridge

```go
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
)

type PrinterAPI interface {
	PrintMessage(string) error
}

type PrinterImpl1 struct{}

func (p *PrinterImpl1) PrintMessage(msg string) error {
	fmt.Printf("%s\n", msg)
	return nil
}

type PrinterImpl2 struct {
	Writer io.Writer
}

func (d *PrinterImpl2) PrintMessage(msg string) error {
	if d.Writer == nil {
		return errors.New("you need to pass an io.Writer to PrinterImpl2")
	}
	fmt.Fprintf(d.Writer, "%s", msg)
	return nil
}

type TestWriter struct {
	Msg string
}

func (t *TestWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 {
		t.Msg = string(p)
		return n, nil
	}
	err = errors.New("content received on writer is empty")
	return
}

type PrinterAbstraction interface {
	Print() error
}
type NormalPrinter struct {
	Msg     string
	Printer PrinterAPI
}

func (c *NormalPrinter) Print() error {
	c.Printer.PrintMessage(c.Msg)
	return nil
}

func main() {
	testWriter := TestWriter{}

	api1 := PrinterImpl1{}
	err := api1.PrintMessage("hello")
	if err != nil {
		log.Fatal(err)
	}

	api2 := PrinterImpl2{
		Writer: &testWriter,
	}
	err = api2.PrintMessage("hello")
	if err != nil {
		log.Fatal(err)
	}

	normal := NormalPrinter{
		Msg:     "hello",
		Printer: &PrinterImpl1{},
	}
	err = normal.Print()
	if err != nil {
		log.Fatal(err)
	}

	normal = NormalPrinter{
		Msg: "hello",
		Printer: &PrinterImpl2{
			Writer: &testWriter,
		},
	}
	err = normal.Print()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Proxy

```go
package main

import (
	"fmt"
	"log"
	"math/rand"
)

type UserFinder interface {
	FindUser(id int32) (User, error)
}

type User struct {
	ID int32
}

type UserList []User

func (u *UserList) FindUser(id int32) (User, error) {
	for i := 0; i < len(*u); i++ {
		if (*u)[i].ID == id {
			return (*u)[i], nil
		}
	}
	return User{}, fmt.Errorf("user %d could not be found", id)
}

func (u *UserList) addUser(newUser User) {
	*u = append(*u, newUser)
}

type UserListProxy struct {
	SomeDatabase           UserList
	StackCache             UserList
	StackCapacity          int
	DidLastSearchUserCache bool
}

func (u *UserListProxy) FindUser(id int32) (User, error) {
	user, err := u.StackCache.FindUser(id)
	if err == nil {
		fmt.Println("returning user from cache")
		u.DidLastSearchUserCache = true
		return user, nil
	}
	user, err = u.SomeDatabase.FindUser(id)
	if err != nil {
		return User{}, err
	}
	u.DidLastSearchUserCache = false
	u.addUserToStack(user)
	return user, nil
}

func (u *UserListProxy) addUserToStack(user User) {
	if len(u.StackCache) >= u.StackCapacity {
		u.StackCache = append(u.StackCache[1:], user)
	} else {
		u.StackCache.addUser(user)
	}
}

func main() {
	someDatabase := UserList{}
	rand.Seed(1)
	for i := 0; i < 1e6; i++ {
		n := rand.Int31()
		someDatabase = append(someDatabase, User{ID: n})
	}
	proxy := UserListProxy{
		SomeDatabase:  someDatabase,
		StackCapacity: 2,
		StackCache:    UserList{},
	}

	knownIDs := [3]int32{
		someDatabase[3].ID,
		someDatabase[4].ID,
		someDatabase[5].ID,
	}

	{
		user, err := proxy.FindUser(knownIDs[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user)
	}
	{
		user, err := proxy.FindUser(knownIDs[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user)
	}
	{
		user, err := proxy.FindUser(someDatabase[rand.Intn(len(someDatabase))-1].ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user)
	}
}
```


## Facade

```go
package main

import (
	"errors"
	"fmt"
	"log"
)

type IngredientAdd interface {
	AddIngredient() (string, error)
}

type PizzaDecorator struct {
	Ingredient IngredientAdd
}

func (p *PizzaDecorator) AddIngredient() (string, error) {
	return "Pizza with the following ingredients:", nil
}

type Meat struct {
	Ingredient IngredientAdd
}

func (m *Meat) AddIngredient() (string, error) {
	if m.Ingredient == nil {
		return "", errors.New("An ingredient is needed in the Ingredient field of the Meat")
	}
	s, err := m.Ingredient.AddIngredient()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", s, "meat"), nil
}

type Onion struct {
	Ingredient IngredientAdd
}

func (o *Onion) AddIngredient() (string, error) {
	if o.Ingredient == nil {
		return "", errors.New("An ingredient is needed in the Ingredient field of the Onion")
	}
	s, err := o.Ingredient.AddIngredient()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", s, "onion"), nil
}

func main() {
	meat := &Onion{&Meat{&PizzaDecorator{}}}
	meatResult, err := meat.AddIngredient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(meatResult)
}
```
