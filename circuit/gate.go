// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package circuit

import (
	"log"
)

// Convenience access

func (u Circuit) IntOrZeroAt(name Name) int {
	i, ok := u.OptionAt(name)
	if !ok {
		return 0
	}
	return i.(int)
}

func (u Circuit) NameAt(name Name) Name {
	return u.At(name).(Name)
}

func (u Circuit) ComplexAt(name Name) complex128 {
	return u.At(name).(complex128)
}

func (u Circuit) FloatAt(name Name) float64 {
	return u.At(name).(float64)
}

func (u Circuit) FloatOrZeroAt(name Name) float64 {
	f, ok := u.OptionAt(name)
	if !ok {
		return 0
	}
	return f.(float64)
}

func (u Circuit) CircuitAt(name Name) Circuit {
	return u.At(name).(Circuit)
}

func (u Circuit) VectorAt(name Name) Vector {
	return u.At(name).(Vector)
}

func (u Circuit) CircuitOptionAt(name Name) (Circuit, bool) {
	v, ok := u.OptionAt(name)
	if ok {
		return v.(Circuit), ok
	}
	return Circuit{}, false
}

func (u Circuit) IntAt(name Name) int {
	return u.At(name).(int)
}

func (u Circuit) IntOptionAt(name Name) (int, bool) {
	v, ok := u.OptionAt(name)
	if ok {
		return v.(int), ok
	}
	return 0, false
}

func (u Circuit) StringAt(name Name) string {
	return u.At(name).(string)
}

func (u Circuit) StringOptionAt(name Name) (string, bool) {
	v, ok := u.OptionAt(name)
	if ok {
		return v.(string), ok
	}
	return "", false
}

func (u Circuit) AddressAt(name Name) Address {
	return u.At(name).(Address)
}

func (u Circuit) AddressOptionAt(name Name) (Address, bool) {
	v, ok := u.OptionAt(name)
	if ok {
		return v.(Address), ok
	}
	return Address{}, false
}

// Series-application methods

func (u Circuit) ReGrow(name Name, value Value) Circuit {
	u.Include(name, value)
	return u
}

func (u Circuit) Grow(name Name, value Value) Circuit {
	if u.Include(name, value) != nil {
		panic("over writing")
	}
	return u
}

func (u Circuit) Refine(name string) Circuit {
	r := New()
	u.Grow(name, r)
	return r
}

func (u Circuit) Abandon(name Name) Circuit {
	u.Exclude(name)
	return u
}

func (u Circuit) Rename(x, y Name) Circuit {
	m := u.Exclude(x)
	if m == nil {
		panic("np")
	}
	if u.Include(y, m) != nil {
		panic("over")
	}
	return u
}

func (u Circuit) Goto(gate ...Name) Value {
	x := u
	for i, g := range gate {
		if i+1 == len(gate) {
			return x.At(g)
		}
		var ok bool
		x, ok = x.CircuitOptionAt(g)
		if !ok {
			log.Fatalf("Address %v points to nothing", Address{gate})
		}
	}
	return x
}

// Assembly

func (u Circuit) Include(name Name, value Value) (before Value) {
	before = u.Gate[name]
	u.Gate[name] = value
	return
}

func (u Circuit) Exclude(name Name) (forgotten Value) {
	forgotten = u.Gate[name]
	delete(u.Gate, name)
	return
}

func (u Circuit) Len() int {
	return len(u.Gate)
}

func (u Circuit) OptionAt(name Name) (Value, bool) {
	v, ok := u.Gate[name]
	return v, ok
}

func (u Circuit) At(name Name) Value {
	return u.Gate[name]
}

const Super = ""
