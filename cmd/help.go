package cmd

import "fmt"

// NOTE: Init help and usage functions
func UsageInit() {
	fmt.Print(`
Default usage: todo init
	Use 'todo init -help' to see more
`)
}
func HelpInit() {
	fmt.Print(`
Help for todo init:
	Available arguments:
		-help  | Show help for todo init
		-h     | Same as '-help'

	Initialize a local todo list in current directory
`)
}

// NOTE: Remove command help and usage functions
func UsageRm() {
	fmt.Print(`
Usage: todo rm [<args>] <title>
	Use 'todo rm --help' to see arguments
`)
}
func HelpRm() {
	fmt.Print(`
Help for todo rm / done:
	Available arguments:
		-help  | Show help for todo rm
		-h     | Same as '--help'
		-all   | Fully remove todo list from current directory
		-a     | Same as '--all'

	Rm and done are the same command with a different name.
	Use 'todo rm <title>' where <title> is the title
	for the list entry to be deleted.
`)
}

// NOTE: Add command help and usage functions
func UsageAdd() {
	fmt.Print(`
Usage: todo add <title>
	Use 'todo add -help' to see more
`)
}
func HelpAdd() {
	fmt.Print(`
Help for todo add:
	Available arguments:
		-help  | Show help for todo add
		-h     | Same as '-help'

	Use 'todo add <title>' where <title> is what you want as
	a title for the new todo entry
`)
}

// NOTE: List command help and usage functions
func UsageList() {
	fmt.Print(`
Usage: todo list [<args>]
	Use 'todo list -help' to see arguments
`)
}
func HelpList() {
	fmt.Print(`
Help for todo list:
	Available arguments:
		-help       | Show help for todo list
		-h          | Same as '-help'
		-all        | List all locations with todo's
		-a          | Same as '-all'
		-pager      | Toggle pagering behavior, by default normal lists
		            | get pagered, locations don't
		-p          | Same as '-pager'

	Show content in a local todo list, or alternatively with '-all'
	show all locations with todo lists
`)
}

// NOTE: Edit command help and usage functions
func UsageEdit() {
	fmt.Print(`
Usage: todo edit [<args>] <title>
	Use 'todo edit -help' to see arguments
`)
}
func HelpEdit() {
	fmt.Print(`
Help for todo edit:
	Available arguments:
		-help  | Show help for todo edit
		-h     | Same as '--help'
		-keep  | New content gets appended to already existing
		       | content instead of overriding it
		-k     | Same as '--keep'

	Edit an existing todo list entry with a given title
`)
}
