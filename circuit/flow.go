// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package circuit

import (
	// "fmt"
	"log"
)

func (u Circuit) Link(x, y Vector) {
	if x.Gate == y.Gate && x.Valve == y.Valve {
		panic("self loop")
	}
	xs, ys := u.valves(x.Gate), u.valves(y.Gate)
	if z, ok := xs[x.Valve]; ok && !Same(z, y) {
		log.Fatalf("%v:%v already connected to %v, not %v", x.Gate, x.Valve, z, y)
		panic("contra")
	}
	if z, ok := ys[y.Valve]; ok && !Same(z, x){
		log.Fatalf("%v:%v already connected to %v, not %v", y.Gate, y.Valve, z, x)
		panic("contra")
	}
	xs[x.Valve], ys[y.Valve] = y, x
}

func (u Circuit) valves(p Name) map[Name]Vector {
	if _, ok := u.Flow[p]; !ok {
		u.Flow[p] = make(map[Name]Vector)
	}
	return u.Flow[p]
}

func (u Circuit) Unlink(x, y Vector) {
	xs, ys := u.Flow[x.Gate], u.Flow[y.Gate]
	delete(xs, x.Valve)
	delete(ys, y.Valve)
	if len(xs) == 0 {
		delete(u.Flow, x.Gate)
	}
	if len(ys) == 0 {
		delete(u.Flow, y.Gate)
	}
}

func (u Circuit) Valves(gate Name) map[Name]Vector {
	return u.Flow[gate]
}

func (u Circuit) ValveNames(gate Name) []Name {
	var r []Name
	for n, _ := range u.Flow[gate] {
		r = append(r, n)
	}
	return r
}

func (u Circuit) View(gate Name) Circuit {
	x := New()
	for vlv, vec := range u.Flow[gate] {
		x.Include(vlv, u.At(vec.Gate))
	}
	return x
}

func (u Circuit) Follow(v Vector) Vector {
	return u.Flow[v.Gate][v.Valve]
}
