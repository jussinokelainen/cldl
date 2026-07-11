package flagger

import (
	"fmt"
	"slices"
)

/*
Function to handle a single argument. If the argument requires a value,
Takes the value at args[i+1], and increments i by 1 to skip the handling of
the value, since it is already used.

Returns an error if and only if there is a flag that requires a
value to be after it, but it is the last value of args.
*/
func handleFlag(flags *Arguments, args []string, currentArg string, i *int) error {
	requiresValue := false
	if slices.Contains(FLAGSET.Valued_flags, currentArg) {
		requiresValue = true
	} else if slices.Contains(FLAGSET.Optional_value, currentArg) {
		requiresValue = true
		if (*i + 1) >= len(args) {
			requiresValue = false
		} else if args[*i+1][0] == '-' {
			requiresValue = false
		}

	}

	if requiresValue {
		if (*i + 1) >= len(args) {
			return fmt.Errorf("value flag as last argument")
		}

		flags.ValueFlags = append(flags.ValueFlags, [2]string{currentArg, args[*i+1]})
		*i++
	} else {
		isValid := false
		if VALIDATE_FLAGS {
			if slices.Contains(FLAGSET.Flags, currentArg) || slices.Contains(FLAGSET.Optional_value, currentArg) {
				isValid = true
				flags.Flags = append(flags.Flags, currentArg)
			}
		} else {
			isValid = true
			flags.Flags = append(flags.Flags, currentArg)
		}
		if !isValid {
			HAS_INVALID = fmt.Errorf("Invalid flag")
		}
	}

	return nil
}
