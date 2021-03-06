// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package handbook

import (
	// "fmt"
	"path"
	"sync"

	. "github.com/gocircuit/escher/image"
	"github.com/gocircuit/escher/be"
	"github.com/gocircuit/escher/see"
	"github.com/gocircuit/escher/faculty"
	"github.com/gocircuit/escher/kit/plumb"
)

func init() {
	ns := faculty.Root.Refine("handbook")
	ns.Grow("Join", Join{}) ??
}

// Join
type Join struct{}

func (Join) Materialize() be.Reflex {
	reflex, _ := be.NewEyeCognizer((&join{}).Cognize, "_", "Head", "Tail")
	return reflex
}

type join struct {
	sync.Mutex
	head *string
	tail *string
}

func (x *join) Cognize(eye *be.Eye, dvalve string, dvalue interface{}) {
	x.Lock()
	defer x.Unlock()
	switch dvalve {
	case "Head":
		head := dvalue.(string)
		x.head = &head
	case "Tail":
		tail := dvalue.(string)
		x.tail = &tail
	default:
		return
	}
	if x.head == nil || x.tail == nil {
		return
	}
	eye.Show("_", path.Join(*x.head, *x.tail))
}
