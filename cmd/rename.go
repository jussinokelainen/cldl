package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Rename a todo entry with given title (change the title)
func RenameTodo(oldTitle string, colors ColorConf) {
	setColorScheme(colors)
	if oldTitle == "" {
		ERROR("Title required")
		return
	}

	_, err := getIfEntryExists(oldTitle)
	if err != nil {
		ERROR("No todo list entry with title:", oldTitle)
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%sEnter new title for todo entry titled %s: \033[0m\n", wipColor, oldTitle)

	fmt.Printf("%s❯ \033[0m", borderColor)
	newTitle, err := reader.ReadString('\n')
	if err != nil {
		ERROR("Error reading input, did text end with a newline?")
		UsageEdit()
		return
	}
	newTitle = strings.TrimSpace(newTitle)
	changeEntryTitle(newTitle, oldTitle)
}

// Change the title of an entry in the database. If checking
// for the title existing is wanted, it should be done before
// calling this function, it does not check it, since the title
// not existing shouldnt break anything in the function
func changeEntryTitle(newTitle string, oldTitle string) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	sqlStatement := `UPDATE todo SET title = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, newTitle, oldTitle)
	if err != nil {
		ERROR("Failed to edit todo content")
		panic(err)
	}

	OK("Successfully changed title", oldTitle, "to", newTitle)
}

func UsageRename() {
	fmt.Print(`Usage: todo rename [-h | --help] <title>
    Use 'todo fix --help' to see more
`)
}

const HelpRename = `Help for todo rename:
    Available arguments:
        --help, -h   | Show this message

    This command is used to change the title
    of and existing todo entry.`
