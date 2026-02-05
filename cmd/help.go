package cmd

import "fmt"

// NOTE: Init help and usage functions
func UsageInit() {
	fmt.Print(`
Usage: todo init [<args>]
	Use todo init --help to see arguments
`)
}
func HelpInit() {
	fmt.Print(`
Help for todo init:
	Available arguments:
		--help | Show help for todo init
		-h     | Same as '--help'

	Initialize a local todo list in current directory
`)
}

// NOTE: Remove command help and usage functions
func UsageRm() {
	fmt.Print(`
Usage: todo rm [<args>] <title>
	Use todo rm --help to see arguments
`)
}
func HelpRm() {
	fmt.Print(`
Help for todo rm:
	Available arguments:
		--help | Show help for todo rm
		-h     | Same as '--help'
		--all  | Fully remove todo list from current dir
		-a     | Same as '--all'

	Use 'todo rm <title>' where <title> is the title
	for the list entry to be deleted
`)
}

// NOTE: Add command help and usage functions
func UsageAdd() {
	fmt.Print(`
Usage: todo add [<args>] <title>
	Use todo add --help to see arguments
`)
}
func HelpAdd() {
	fmt.Print(`
Help for todo add:
	Available arguments:
		--help | Show help for todo add
		-h     | Same as '--help'

	Use 'todo add <title>' where <title> is the title
	for the new todo entry
`)
}

// NOTE: List command help and usage functions
func UsageList() {
	fmt.Print(`
Usage: todo list [<args>]
	Use todo list --help to see arguments
`)
}
func HelpList() {
	fmt.Print(`
Help for todo list:
	Available arguments:
		--help      | Show help for todo list
		-h          | Same as '--help'
		--all       | List all locations with todo's
		-a          | Same as '--all'
		--pager     | ONLY with --all, sends the list to a pager
		--no-pager  | Don't use a pager for printing the local list

	Show all todo list entries, or all todo lists
`)
}

// NOTE: Edit command help and usage functions
func UsageEdit() {
	fmt.Print(`
Usage: todo edit [<args>] <title>
	Use todo edit --help to see arguments
`)
}
func HelpEdit() {
	fmt.Print(`
Help for todo edit:
	Available arguments:
		--help | Show help for todo edit
		-h     | Same as '--help'
		--keep | New content gets appended to already existing
		       | content instead of overriding it
		-k     | Same as '--keep'

	Edit a todo list entry that already exists
	with a title given as argument
`)
}
