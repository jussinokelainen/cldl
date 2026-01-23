package main

import "fmt"

// NOTE: Main help and usage functions
func mainUsage() {
	fmt.Print(`
Usage: todo <COMMAND> [<args>]
	Use todo --help to see available commands and arguments
`)
}
func mainHelp() {
	fmt.Print(`
Help for todo:
	Available commands:
		--help		   | Show help message
		-h			   | Same as '--help'
		init		   | Create new todo in current dir
		list		   | List all todo list entries
		add			   | Add new entry into todo list
		rm <title>	   | Remove todo list entry or entire list, see 'todo rm --help'
		remove <title> | Same as 'rm'
		done <title>   | Same as 'rm'

	Todo application that creates local per-directory todo-lists with sqlite
`)
}

// NOTE: Init help and usage functions
func usageInit() {
	fmt.Print(`
Usage: todo init [<args>]
	Use todo init --help to see arguments
`)
}
func helpInit() {
	fmt.Print(`
Help for todo init:
	Available arguments:
		--help | Show help for todo init
		-h	   | Same as '--help'

	Initialize a local todo list in current directory
`)
}

// NOTE: Remove command help and usage functions
func usageRm() {
	fmt.Print(`
Usage: todo rm [<args>] <title>
	Use todo rm --help to see arguments
`)
}
func helpRm() {
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
func usageAdd() {
	fmt.Print(`
Usage: todo add [<args>] <title>
	Use todo add --help to see arguments
`)
}
func helpAdd() {
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
func usageList() {
	fmt.Print(`
Usage: todo list [<args>]
	Use todo list --help to see arguments
`)
}
func helpList() {
	fmt.Print(`
Help for todo list:
	Available arguments:
		--help | Show help for todo list
		-h     | Same as '--help'
		--all  | List all locations with todo's
		-a     | Same as '--all'

	Show all todo list entries, or all todo lists
`)
}
