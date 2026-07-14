package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Rename a todo entry with given title (change the title)
func Rename_todo(oldTitle string, colors ColorConf) {
	set_color_scheme(colors)
	if oldTitle == "" {
		ERROR("Title required")
		return
	}

	_, err := get_content_if_entry_exists(oldTitle)
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
		Usage_edit()
		return
	}
	newTitle = strings.TrimSpace(newTitle)
	change_entry_title(newTitle, oldTitle)
}

// Change the title of an entry in the database. If checking
// for the title existing is wanted, it should be done before
// calling this function, it does not check it, since the title
// not existing shouldnt break anything in the function
func change_entry_title(newTitle string, oldTitle string) {
	todoDB := open_todo_db()
	defer todoDB.Close()

	sqlStatement := `UPDATE todo SET title = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, newTitle, oldTitle)
	if err != nil {
		ERROR("Failed to edit todo content")
		panic(err)
	}

	OK("Successfully changed title", oldTitle, "to", newTitle)
}

func Usage_rename() {
	fmt.Print(`Usage: cldl rename [-h | --help] <title>
    Use 'cldl fix --help' to see more
`)
}

const HelpRename = `Help for cldl rename:
    Available arguments:
        --help, -h   | Show this message

    This command is used to change the title
    of and existing todo entry.`
