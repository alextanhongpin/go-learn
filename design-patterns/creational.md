# Singleton

```go
package singleton

import "sync"

var instance *singleton

var once sync.Once

type singleton struct{}

func GetInstance() *singleton {
	once.Do(func() {
		instance = new(singleton)
	})
	return instance
}
```

# Builder

```go
package main

import (
	"fmt"
)

type BuildProcess interface {
	SetWheels() BuildProcess
	SetSeats() BuildProcess
	SetStructure() BuildProcess
	Build() VehicleProduct
}

type ManufacturingDirector struct {
	builder BuildProcess
}

func (f *ManufacturingDirector) Construct() {
	f.builder.SetSeats().SetStructure().SetWheels()
}

func (f *ManufacturingDirector) SetBuilder(b BuildProcess) {
	f.builder = b
}

type VehicleProduct struct {
	Wheels    int
	Seats     int
	Structure string
}

type CarBuilder struct {
	v VehicleProduct
}

func (c *CarBuilder) SetWheels() BuildProcess {
	c.v.Wheels = 4
	return c
}

func (c *CarBuilder) SetSeats() BuildProcess {
	c.v.Seats = 5
	return c
}

func (c *CarBuilder) SetStructure() BuildProcess {
	c.v.Structure = "Car"
	return c
}

func (c *CarBuilder) Build() VehicleProduct {
	return c.v
}

type BikeBuilder struct {
	v VehicleProduct
}

func (b *BikeBuilder) SetWheels() BuildProcess {
	b.v.Wheels = 2
	return b
}

func (b *BikeBuilder) SetSeats() BuildProcess {
	b.v.Seats = 2
	return b
}

func (b *BikeBuilder) SetStructure() BuildProcess {
	b.v.Structure = "Motorbike"
	return b
}

func (b *BikeBuilder) Build() VehicleProduct {
	return b.v
}

func main() {
	manufacturingComplex := ManufacturingDirector{}
	carBuilder := new(CarBuilder)
	manufacturingComplex.SetBuilder(carBuilder)
	manufacturingComplex.Construct()
	car := carBuilder.Build()
	fmt.Printf("%#v\n", car)

	bikeBuilder := new(BikeBuilder)
	manufacturingComplex.SetBuilder(bikeBuilder)
	manufacturingComplex.Construct()
	bike := bikeBuilder.Build()
	fmt.Printf("%#v\n", bike)
}
```

# Factory
```go
package main

import (
	"fmt"
	"log"
)

type PaymentMethod interface {
	Pay(amount float32) string
}

const (
	Cash      = 1
	DebitCard = 2
)

func GetPaymentMethod(m int) (PaymentMethod, error) {
	switch m {
	case Cash:
		return new(CashPM), nil
	case DebitCard:
		return new(DebitCardPM), nil
	default:
	}
	return nil, fmt.Errorf("payment method %d not recognized", m)
}

type CashPM struct{}
type DebitCardPM struct{}

func (c *CashPM) Pay(amount float32) string {
	return fmt.Sprintf("paid %0.2f using cash", amount)
}

func (c *DebitCardPM) Pay(amount float32) string {
	return fmt.Sprintf("paid %0.2f using debit card", amount)
}

func main() {
	payment, err := GetPaymentMethod(Cash)
	if err != nil {
		log.Fatal(err)
	}
	msg := payment.Pay(10.30)
	fmt.Println(msg)
}
```

# Abstract Factory

```go
package main

import (
	"fmt"
	"log"
)

type Vehicle interface {
	NumWheels() int
	NumSeats() int
}

type Car interface {
	NumDoors() int
}

type Motorbike interface {
	GetMotorbikeType() int
}

type VehicleFactory interface {
	Build(v int) (Vehicle, error)
}

const (
	LuxuryCarType = 1
	FamilyCarType = 2

	CarFactoryType = 1
)

type CarFactory struct{}

func (c *CarFactory) Build(v int) (Vehicle, error) {
	switch v {
	case LuxuryCarType:
		return new(LuxuryCar), nil
	case FamilyCarType:
		return new(FamilyCar), nil
	default:
		return nil, fmt.Errorf("vehicle of type %d not recognized", v)
	}
}

type LuxuryCar struct {
}

func (l *LuxuryCar) NumDoors() int {
	return 4
}
func (l *LuxuryCar) NumWheels() int {
	return 4
}

func (l *LuxuryCar) NumSeats() int {
	return 5
}

type FamilyCar struct {
}

func (f *FamilyCar) NumDoors() int {
	return 5
}
func (f *FamilyCar) NumWheels() int {
	return 4
}

func (f *FamilyCar) NumSeats() int {
	return 5
}

func BuildFactory(f int) (VehicleFactory, error) {
	switch f {
	case CarFactoryType:
		return new(CarFactory), nil
	default:
		return nil, fmt.Errorf("factory with id %d not recognized", f)
	}

}

func main() {
	fmt.Println("Hello, playground")
	carF, err := BuildFactory(CarFactoryType)
	if err != nil {
		log.Fatal(err)
	}
	carVehicle, err := carF.Build(LuxuryCarType)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(carVehicle.NumWheels())
	car, ok := carVehicle.(*LuxuryCar)
	if ok {
		fmt.Println(car.NumDoors())
	}
}
```

# Prototype

```go
package main

import (
	"errors"
	"fmt"
	"log"
)

type ShirtCloner interface {
	GetClone(s int) (ItemInfoGetter, error)
}

const (
	White = 1
	Black = 2
	Blue  = 3
)

type ShirtsCache struct {
}

func (s *ShirtsCache) GetClone(i int) (ItemInfoGetter, error) {
	switch i {
	case White:
		newItem := *whitePrototype
		return &newItem, nil
	default:
		return nil, errors.New("shirt model not recognized")
	}
}

type ItemInfoGetter interface {
	GetInfo() string
}

type ShirtColor byte
type Shirt struct {
	Price float32
	SKU   string
	Color ShirtColor
}

func (s *Shirt) GetInfo() string {
	return ""
}
func GetShirtsCloner() ShirtCloner {
	return &ShirtsCache{}
}

var whitePrototype *Shirt = &Shirt{
	Price: 15.00,
	SKU:   "empty",
	Color: White,
}

var blackPrototype *Shirt = &Shirt{
	Price: 16.00,
	SKU:   "empty",
	Color: Black,
}

var bluePrototype *Shirt = &Shirt{
	Price: 17.00,
	SKU:   "empty",
	Color: Blue,
}

func (i *Shirt) GetPrice() float32 {
	return i.Price
}

func main() {
	shirtCache := GetShirtsCloner()
	if shirtCache == nil {
		log.Fatal("cache is nil")
	}
	item1, err := shirtCache.GetClone(White)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%#v - %p\n", item1, item1)
	item2, err := shirtCache.GetClone(White)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v - %p\n", item2, item2)
}
```
