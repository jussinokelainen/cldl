package flagger

import (
	"fmt"
)

var (
	FLAGSET        Flagset
	HAS_INVALID    error
	VALIDATE_FLAGS bool
)

/*
Struct to pass valid flags to the parser. If it is desired to set all flags as
valid in a certain type of flag, the value can be set as nil.

	Flags          | Normal flags (bool i guess)
	Valued_flags   | Flags that take an argument as their value (like --depth 1)
	Optional_value | Flags that can take an argument, but can also be used without

Optionally valued flags will try to get a value, unless the flag is the last argument,
or the next argument is another flag (starts with '-')
*/
type Flagset struct {
	Flags          []string
	Valued_flags   []string
	Optional_value []string
}

type Arguments struct {
	Flags      []string
	ValueFlags [][2]string
	NormalStr  []string
}

/*
Parse arguments, (intended to be given something like os.Args[1:]), returns
them as a struct with separate slices for non-arguments, normal flags and
arguments that take values after them (e.g. '--length 1').

Takes as a parameter a Flagset struct, which contains all valid flags or nil,
to accept flags in the corresponding flag type

Returns an error if arguments contain invalid flags or a flag that requires a
value is the last argument
*/
func ParseFlags(args []string, validFlags Flagset) (Arguments, error) {
	FLAGSET = validFlags
	HAS_INVALID = nil

	if validFlags.Flags == nil {
		VALIDATE_FLAGS = false
	} else {
		VALIDATE_FLAGS = true
	}

	var flags Arguments
	flags.Flags = nil
	flags.ValueFlags = nil
	flags.NormalStr = nil

	for i := 0; i < len(args); i++ {
		currentArg := args[i]

		if currentArg[0] == '-' {
			if len(currentArg) == 1 {
				HAS_INVALID = fmt.Errorf("Empty flag")
			} else if currentArg[1] == '-' {
				if len(currentArg) == 2 {
					HAS_INVALID = fmt.Errorf("Empty flag")
				} else {
					// --[argument] style arg
					err := handleFlag(&flags, args, currentArg[2:], &i)
					if err != nil {
						HAS_INVALID = fmt.Errorf("Value flag cannot be last argument")
					}
				}
			} else {
				// -[arg][arg][arg] style arg
				for _, c := range currentArg[1:] {
					err := handleFlag(&flags, args, string(c), &i)
					if err != nil {
						HAS_INVALID = fmt.Errorf("Value flag cannot be last argument")
					}
				}
			}

		} else {
			flags.NormalStr = append(flags.NormalStr, currentArg)
		}
	}

	return flags, HAS_INVALID
}
