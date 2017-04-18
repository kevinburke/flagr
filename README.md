# flagr

This exists because Go won't parse flag arguments ("-o" or "--opt") that come
after positional arguments:

```go
f := flag.NewFlagSet("foo", flag.ExitOnError)
opt := f.String("opt", "", "an option")
f.Parse([]string{"regular-argument", "--opt", "theoption"})
fmt.Println(f.Args())
fmt.Println(*opt) // ""
```

Sometimes this does the right thing - "go" takes different flags than "go build"
for example and you don't want to parse "go" flags placed after "build". However
sometimes you are at the end of the parse chain and anything that's left should
be counted as an argument.

This does the dumbest possible thing which is to rearrange the flags so any
options come first. Anything after `--` is ignored. A single `-` counts as a
positional (i.e. non-flag) argument.

```go
f := flag.NewFlagSet("test", flag.ContinueOnError)
f.String("blah", "", "Blah")
f.Bool("baz", false, "Baz")
args := Rearrange(f, []string{"foo", "two", "--blah", "bar", "another", "--str=value", "--baz", "--", "--last", "option"})
fmt.Println(args)
// Output: [--blah bar --str=value --baz foo two another -- --last option]
```
