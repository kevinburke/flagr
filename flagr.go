// Package flagr rearranges flags so the options come first.
package flagr

import (
	"flag"
	"fmt"
)

type Booler interface {
	IsBoolFlag() bool
}

// insertAndDelete pops the value at j and inserts it at i, where i < j.
func insertAndDelete(args []string, i, j int) {
	if i >= j {
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
	if len(args) == 0 {
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
