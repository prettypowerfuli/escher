// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package shell

import (
	"bufio"
	"fmt"
	"io"
	"path"
	"strings"
	// "time"

	// . "github.com/gocircuit/escher/faculty"
	. "github.com/gocircuit/escher/circuit"
	// . "github.com/gocircuit/escher/be"
	// . "github.com/gocircuit/escher/fs"
	. "github.com/gocircuit/escher/memory"
)

type Shell struct {
	memory *Memory
	in io.Reader
	out io.WriteCloser
	err io.WriteCloser
	scan *bufio.Scanner
	foci map[string]*Focus
	focus string
}

type Focus struct {
	Name string
	Address Address
}

func NewShell(in io.Reader, out, err io.WriteCloser, memory *Memory) *Shell {
	sh := &Shell{
		memory: memory,
		in: in,
		out: out,
		err: err,
		scan: bufio.NewScanner(in),
		foci: make(map[string]*Focus),
		focus: "a",
	}
	for _, c := range "abcdefghijklmnopqrstuvwxyz" {
		sh.foci[string(c)] = &Focus{
			Name: string(c),
			Address: NewAddress([]Name{}),
		}
	}
	return sh
}

func (sh *Shell) Loop() {
	sh.ShowPrompt()
	for sh.scan.Scan() {
		words := split(sh.scan.Text())
		if len(words) == 0 {
			sh.ShowPrompt()
			continue
		}
		switch strings.ToLower(words[0]) {
		case "help", "h", "?":
			sh.ShowHelp()
		case "ls":
			sh.ShowCircuit(words[1:])
		case "mk", "incl", "excl":
			sh.Include(words[1:])
		}
		sh.ShowPrompt()
	}
}

func (sh *Shell) ShowPrompt() {
	fmt.Fprintf(sh.err, " ")
}

func (sh *Shell) ShowCircuit(w []string) {
	if len(w) == 0 {
		fmt.Fprintf(sh.err, "%v\n", sh.memory.View())
		return
	}
	// XXX
	fmt.Fprintf(sh.err, "!@#$\n")
}

func (sh *Shell) Include(w []string) {
	switch {
	case len(w) == 0:
		fmt.Fprintf(sh.err, "Make command needs arguments\n")
	case len(w) == 1:
		dir, file := path.Split(w[0])
		if dir != "" {
			fmt.Fprintf(sh.err, "Make argument cannot be a path\n")
			return
		}
		sh.memory.Include(file, New())
	}
}

func (sh *Shell) ShowFocus() {
	fmt.Fprintf(sh.err, "%s:%v\n", sh.focus, "/" + strings.Join(sh.foci[sh.focus].Address.Strings(), "/"))
}

func (sh *Shell) ShowHelp() {
	const help = `
help		Show this screen
cd		Move current focus
ls		Show the current circuit
mk		Add a sub-circuit
fwd, f	Move forward in history
bwd, b	Move backward in history
jump		Change the current focus
`
	fmt.Fprintf(sh.err, "%s\n", help)
}

// func glob(src string) (addr Address, recurse bool) {
// 	??
// }

func split(src string) (r []string) {
	var s int
	for i, b := range src {
		if b != ' ' && b != '\t' {
			continue
		}
		switch {
		case i > s:
			r = append(r, src[s:i])
			s = i + 1
		case i == s:
			s++
		}
	}
	if len(src) > s {
		r = append(r, src[s:])
	}
	return
}