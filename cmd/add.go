package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Adds a new todo with a title given as an argument, if title is not duplicate
func AddTodo(title string) {
	if !TodoExists() {
		fmt.Print("[\033[35m INFO \033[0m] No todo currently exists in this directory.\n")
		answer := askIfInit()
		if answer {
			InitTodo()
			fmt.Print("\n")
		} else {
			return
		}
	}

	// Check whether an entry already exists with the same title
	_, exists := getIfEntryExists(title)
	if exists == nil {
		errout("Failed to add new todo. Please select a unique title.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	todoDB := openTodoDB()
	defer todoDB.Close()

	time := time.Now().Unix()

	fmt.Printf("\033[36mEnter contents for new todo titled %s: \033[0m\n", title)
	fmt.Print("\033[35m❯ \033[0m")
	content, err := reader.ReadString('\n')
	if err != nil {
		errout("Error reading input, did text end with a newline")
		UsageAdd()
		return
	}
	content = strings.TrimSpace(content)

	sqlStatement := `INSERT INTO todo(title, content, time) VALUES($1, $2, $3);`
	_, err = todoDB.Exec(sqlStatement, title, content, time)
	if err != nil {
		errout("Error adding new todo, executing database query failed")
		panic(err)
	}
	ok("Successfully added new todo " + title)
}

func askIfInit() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to initialize a new one? [y/n]: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		errout("Error reading input")
		return askIfInit()
	}
	answer = strings.TrimSpace(answer)
	switch answer {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Print("Invalid answer, try again.\n")
		return askIfInit()
	}
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
	a title for the new todo entry.

	Inputting content will be finished by inputting a newline,
	(most likely by pressing 'enter'), Errors may otherwise
`)
}
