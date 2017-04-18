// Package flagr rearranges flags so the options come first.
/*
This exists because Go won't parse flag arguments ("-o" or "--opt") that come
after positional arguments:

	f := flag.NewFlagSet("foo", flag.ExitOnError)
	opt := f.String("opt", "", "an option")
	f.Parse([]string{"regular-argument", "--opt", "theoption"})
	fmt.Println(f.Args())
	fmt.Println(*opt) // ""

Sometimes this does the right thing - "go" takes different flags than "go build"
for example and you don't want to parse "go" flags placed after "build". However
sometimes you are at the end of the parse chain and anything that's left should
be counted as an argument.

This does the dumbest possible thing which is to rearrange the flags so any
options come first. Anything after `--` is ignored. A single `-` counts as a
positional (i.e. non-flag) argument.
*/
package flagr

import (
	"flag"
	"fmt"
)

// Booler tells you whether a flag.Value is a BoolFlag.
type Booler interface {
	// IsBoolFlag returns true if the flag.Value is a boolean flag.
	IsBoolFlag() bool
}

// insertAndDelete pops the value at j and inserts it at i, where i < j.
func insertAndDelete(args []string, i, j int) {
	if i == j {
		return // nothing to do
	}
	if i > j {
		panic(fmt.Sprint("i > j", i, j))
	}
	val := args[j]
	// delete
	args = append(args[:j], args[j+1:]...)

	// insert at i
	args = append(args, "")
	copy(args[i+1:], args[i:])
	args[i] = val
}

// Rearrange rearranges args so any flag arguments (arguments beginning with "-"
// or "--") are placed before the non-flag arguments. Arguments after the
// terminator "--" are returned as is (they are not moved).
func Rearrange(set *flag.FlagSet, args []string) []string {
	if len(args) == 0 || len(args) == 1 {
		return args
	}
	idx := len(args) - 1
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) == 0 || arg[0] != '-' || len(arg) == 1 {
			// non flag argument
			if idx > i {
				idx = i
			}
			continue
		}
		numMinuses := 1
		if arg[1] == '-' {
			numMinuses++
			if len(arg) == 2 {
				return args
			}
		}
		name := arg[numMinuses:]
		hasValue := false
		for j := 1; j < len(name); j++ { // equals cannot be first
			if name[j] == '=' {
				hasValue = true
				name = name[0:j]
				break
			}
		}
		flg := set.Lookup(name)
		if flg == nil {
			// swap it, continue
			insertAndDelete(args, idx, i)
			idx++
			continue
		}
		if fv, ok := flg.Value.(Booler); ok && fv.IsBoolFlag() {
			// no next argument
			insertAndDelete(args, idx, i)
			idx++
			continue
		}
		if hasValue {
			// val attached via equals
			insertAndDelete(args, idx, i)
			idx++
			continue
		}
		// swap the next argument through as well. for example (3, 4) should get
		// swapped to (1, 2)
		insertAndDelete(args, idx, i)
		if i == len(args)-1 {
			idx++
			continue
		}
		insertAndDelete(args, idx+1, i+1)
		idx += 2
		i++
	}
	return args
}
