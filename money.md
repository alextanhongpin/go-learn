# Money

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math"
	"math/big"
)

func main() {
	m := NewMoney(102)
	fmt.Println(m.AllocateEqual(3))
	fmt.Println(m.AllocateRatio(1, 2, 1)) // NOTE: The order matters - the last one will always be +/- 1 Rp.
}

type Money struct {
	value int64
}

func NewMoney(value int64) *Money {
	return &Money{value}
}

func (m Money) AllocateEqual(n int) []int64 {
	res := make([]int64, n)

	if m.value < int64(n) {
		panic(fmt.Errorf("insufficient funds to allocate %d between %d", m.value, n))
	}
	r := big.NewRat(m.value, int64(n))
	f64, exact := r.Float64()
	if exact {
		for i := 0; i < n; i++ {
			res[i] = int64(f64)
		}
		return res
	}

	i64 := int64(math.Round(f64))
	for i := 0; i < n-1; i++ {
		res[i] = i64
	}
	res[n-1] = m.value - (i64 * int64(n-1))
	return res
}

func (m Money) AllocateRatio(ratios ...int64) []int64 {
	var totalRatio int64
	for _, ratio := range ratios {
		totalRatio += ratio
	}

	res := make([]int64, len(ratios))
	lastAmount := m.value
	for i := 0; i < len(ratios)-1; i++ {
		r := big.NewRat(ratios[i], totalRatio)
		r = r.Mul(r, big.NewRat(m.value, 1))
		f64, _ := r.Float64()
		i64 := int64(math.Round(f64))
		res[i] = i64
		lastAmount -= i64
	}
	res[len(ratios)-1] = lastAmount
	return res
}
```


## With Base Unit

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math"
	"math/big"
)

func main() {
	m := NewMoney(1300, 100)
	fmt.Println(m.AllocateEqual(3))
	fmt.Println(m.AllocateRatio(1, 2, 3)) // NOTE: The order matters - the last one will always be +/- 1 Rp.
}

type Money struct {
	value int64
	unit  int64
}

func NewMoney(value, unit int64) *Money {
	return &Money{value, unit}
}

func (m Money) AllocateEqual(n int) []int64 {
	res := make([]int64, n)

	if m.value < int64(n) {
		panic(fmt.Errorf("insufficient funds to allocate %d between %d", m.value, n))
	}

	r := big.NewRat(m.value, m.unit)
	r = r.Mul(r, big.NewRat(1, int64(n)))
	f64, exact := r.Float64()
	if exact {
		for i := 0; i < n; i++ {
			res[i] = int64(f64)
		}
		return res
	}

	lastAmount := m.value
	i64 := int64(math.Round(f64))
	for i := 0; i < n-1; i++ {
		res[i] = i64 * m.unit
		lastAmount -= i64 * m.unit
	}
	res[n-1] = lastAmount
	return res
}

func (m Money) AllocateRatio(ratios ...int64) []int64 {
	var totalRatio int64
	for _, ratio := range ratios {
		totalRatio += ratio
	}

	res := make([]int64, len(ratios))
	lastAmount := m.value
	for i := 0; i < len(ratios)-1; i++ {
		r := big.NewRat(ratios[i], totalRatio)
		r = r.Mul(r, big.NewRat(m.value, 1))
		r = r.Mul(r, big.NewRat(1, int64(m.unit)))
		f64, _ := r.Float64()
		i64 := int64(math.Round(f64))
		res[i] = i64 * m.unit
		lastAmount -= i64 * m.unit
	}
	res[len(ratios)-1] = lastAmount
	return res
}
```
