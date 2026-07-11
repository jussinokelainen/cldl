package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AddInfo struct {
	Priority      int
	Tag           string
	Empty_content bool
	File_path     string
	File_line     int
}

// Adds a new todo with a title given as an argument, if title is not duplicate
func AddTodo(title string, conf AddConf, data AddInfo) {
	if !TodoExists() {
		var initNew bool
		if conf.Auto_init {
			initNew = true
		} else {
			INFO("Not todo currently exists in this directory")
			initNew = askIfInit()
		}
		if initNew {
			InitTodo()
			fmt.Print("\n")
		} else {
			return
		}
	}

	// Check for empty title and whether an entry already exists with the same title
	if title == "" {
		INFO("Title is required")
		UsageAdd()
		return
	}
	_, exists := get_content_if_entry_exists(title)
	if exists == nil {
		ERROR("Failed to add new todo. Please select a unique title.")
		return
	}

	time := time.Now().Unix()

	todoDB := openTodoDB()
	defer todoDB.Close()

	var content string
	if data.Empty_content {
		content = "[ EMPTY ]"
	} else {
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("\033[36mEnter contents for new todo titled %s: \033[0m\n", title)
		fmt.Print("\033[35m❯ \033[0m")
		var err error
		content, err = reader.ReadString('\n')
		if err != nil {
			ERROR("Error reading input, did text end with a newline")
			UsageAdd()
			return
		}
		content = strings.TrimSpace(content)
		if content == "" {
			content = "[ EMPTY ]"
		}
	}

	if conf.Ask_priority {
		data.Priority = askPriority()
	}

	sqlStatement := `INSERT INTO todo(title, content, time, priority, tag, file, line) VALUES($1, $2, $3, $4, $5, $6, $7);`
	_, err := todoDB.Exec(sqlStatement, title, content, time, data.Priority, data.Tag, data.File_path, data.File_line)
	if err != nil {
		ERROR("Error adding new todo, executing database query failed")
		panic(err)
	}

	ext := filepath.Ext(data.File_path)
	comment_pf, ok := lineComments[ext]
	if !ok {
		comment_pf = "#"
	}

	comment_line := fmt.Sprint(comment_pf, " TODO-ENTRY: title: ", title, ", priority: ", data.Priority, ", tag: ", data.Tag)
	insert_line(data.File_path, (data.File_line - 1), comment_line)

	OK("Successfully added new todo " + title)
}

func askPriority() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\033[35mEnter priority to be set for this entry\033[0m: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		ERROR("Error reading input")
		return askPriority()
	}
	answer = strings.TrimSpace(answer)

	answerInt, err := strconv.Atoi(answer)
	if err != nil {
		ERROR("Input must be an integer")
		return askPriority()
	}
	return answerInt
}

func askIfInit() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to initialize a new one? [y/n]: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		ERROR("Error reading input")
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
	fmt.Print(`Usage: todo add [-h | --help] [--auto-init] [-t | --tag] [-p | --priority] <title>
    Use 'todo add --help' to see more
`)
}

const HelpAdd = `Help for todo add:
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
        --tag, -t      | Set the tag for a new entry
        --empty, -e    | Adds todo without asking for content
        --file, -f     | Specify the file connected to the entry
        --line, -l     | Specify the line in the file

    Use 'todo add <title>' where <title> is what you want as a title for the
    new todo entry. Titles can be entered with spaces in them without having
    to use quotes

    Inputting content will be finished by inputting a newline,
    (most likely by pressing 'enter'), Errors may otherwise occur.
    Content can be empty by inputting nothing, or using the --empty argument.

    Config option 'auto_init' Determines whether a new local database is
    created automatically or asked before doing it`
