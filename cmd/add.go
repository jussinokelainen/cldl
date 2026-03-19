package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Adds a new todo with a title given as an argument, if title is not duplicate
func AddTodo(title string, conf AddConf, priority int) {
	if !TodoExists() {
		var answer bool
		if conf.Auto_init {
			answer = true
		} else {
			info("Not todo currently exists in this directory")
			answer = askIfInit()
		}
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
	if content == "" {
		content = "[ EMPTY ]"
	}

	if conf.Ask_priority {
		priority = askPriority()
	}

	sqlStatement := `INSERT INTO todo(title, content, time, priority, tag) VALUES($1, $2, $3, $4, $5);`
	_, err = todoDB.Exec(sqlStatement, title, content, time, priority, "test")
	if err != nil {
		errout("Error adding new todo, executing database query failed")
		panic(err)
	}
	ok("Successfully added new todo " + title)
}

func askPriority() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\033[35mEnter priority to be set for this entry\033[0m: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		errout("Error reading input")
		return askPriority()
	}
	answer = strings.TrimSpace(answer)

	answerInt, err := strconv.Atoi(answer)
	if err != nil {
		errout("Input must be an integer")
		return askPriority()
	}
	return answerInt
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
Usage: todo add [-h | --help] [--auto-init] [-p | --priority] <title>
    Use 'todo add --help' to see more
`)
}
func HelpAdd() {
	fmt.Print(`
Help for todo add:
    Available arguments:
        --help, -h     | Show this message
        --auto-init    | Automatically initialize a new todo when adding
                       | an entry and a list doesn't exist yet.
                       | Can be used without a value to set auto-init
                       | to true, or with a true/false to set the value
                       | (If used without giving the value, it must be after
                       | the title of the entry, otherwise the first word of
                       | the title is interpreted as the value)
        --priority, -p | Specify the priority that will be set for this entry
                       | regardless of default_priority, and it will not be
                       | asked later regardless of ask_priority

    Use 'todo add <title>' where <title> is what you want as a title for the
    new todo entry. Titles can be entered with spaces in them without having
    to use quotes

    Inputting content will be finished by inputting a newline,
    (most likely by pressing 'enter'), Errors may otherwise occur

    Config option 'auto_init' Determines whether a new local database is
    created automatically or asked before doing it
`)
}
