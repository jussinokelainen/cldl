package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

func EditTodo(title string, keep bool) {
	if title == "" {
		errout("Title required")
		return
	}

	content, err := getIfEntryExists(title)
	if err != nil {
		errout("No todo list entry found with title " + title)
		return
	}

	content = wordwrap.WrapString(content, uint(maxWidth))
	fmt.Printf("\033[36mOld content for %s:\n", title)
	fmt.Printf("\033[32m%s\033[0m\n\n", content)

	reader := bufio.NewReader(os.Stdin)
	if keep {
		fmt.Printf("\033[36mEnter content to be added into todo titled %s: \033[0m\n", title)
	} else {
		fmt.Printf("\033[36mEnter new content for todo titled %s: \033[0m\n", title)
	}

	fmt.Print("\033[35m❯ \033[0m")
	newContent, err := reader.ReadString('\n')
	if err != nil {
		errout("Error reading input, did text end with a newline?")
		UsageEdit()
		return
	}
	newContent = strings.TrimSpace(newContent)
	if keep {
		newContent = content + "\n\nNew edit:\n" + newContent
	}
	changeEntryContent(newContent, title)
}

// Send the new content of an entry into the database
func changeEntryContent(newContent string, title string) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	sqlStatement := `UPDATE todo SET content = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, newContent, title)
	if err != nil {
		errout("Failed to edit todo content")
		panic(err)
	}

	fmt.Print("\n")
	ok("Successfully changed content for " + title)

}

// NOTE: Edit command help and usage functions
func UsageEdit() {
	fmt.Print(`
Usage: todo edit [-h | --help] [-k | --keep] <title>
    Use 'todo edit --help' to see more
`)
}
func HelpEdit() {
	fmt.Print(`
Help for todo edit:
    Available arguments:
        --help, -h  | Show this message
        --keep, -k  | Toggles the behavior of keeping on edit

    Edit an existing todo list entry with a given title. Can be configured
    to either keep the existing content or override it by default

    Same text inputting rules apply for editing as adding a new entry,
    Check 'todo add --help'.
`)
}
