package flagr

import (
	"flag"
	"fmt"
	"reflect"
	"testing"
)

func TestRearrange(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.String("blah", "", "Blah")
	f.Bool("baz", false, "Baz")
	args := Rearrange(f, []string{"foo", "two", "--blah", "bar", "another", "--str=value", "--baz", "--", "--last", "option"})
	want := []string{"--blah", "bar", "--str=value", "--baz", "foo", "two", "another", "--", "--last", "option"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("got %q\nwant %q", args, want)
	}
}

func TestOne(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	args := Rearrange(f, []string{"-h"})
	want := []string{"-h"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("got %q\nwant %q", args, want)
	}
}

func Example() {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.String("blah", "", "Blah")
	f.Bool("baz", false, "Baz")
	args := Rearrange(f, []string{"foo", "two", "--blah", "bar", "another", "--str=value", "--baz", "--", "--last", "option"})
	fmt.Println(args)
	// Output: [--blah bar --str=value --baz foo two another -- --last option]
}
