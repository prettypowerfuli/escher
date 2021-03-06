// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package time

import (
	"time"

	"github.com/gocircuit/escher/faculty"
	"github.com/gocircuit/escher/be"
	. "github.com/gocircuit/escher/circuit"
)

func init() {
	faculty.Register("time.Ticker", Ticker{})
	faculty.Register("time.Delay", Delay{})
}

// Delay…
type Delay struct{}

func (Delay) Materialize() (be.Reflex, Value) {
	reflex, eye := be.NewEye("X", "Y", "Duration")
	go func() {
		ds := make(chan time.Duration, 2)
		dur := ds
		xy, yx := make(chan interface{}, 1), make(chan interface{}, 1)
		go func() {
			d := <-dur
			for {
				v := <-xy
				time.Sleep(d)
				eye.Show("Y", v)
			}
		}()
		go func() {
			d := <-dur
			for {
				v := <-yx
				time.Sleep(d)
				eye.Show("X", v)
			}
		}()
		for {
			valve, value := eye.See()
			switch valve {
			case "X":
				xy <- value
			case "Y":
				yx <- value
			case "Duration":
				if ds == nil {
					break
				}
				d := time.Duration(value.(int))
				ds <- d
				ds <- d
				close(ds)
				ds = nil
			}
		}
	}()
	return reflex, Delay{}
}
