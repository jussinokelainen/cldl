package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

func Edit_todo(title string, conf EditConf, colors ColorConf) {
	set_color_scheme(colors)
	if title == "" {
		ERROR("Title required")
		return
	}

	content, err := get_content_if_entry_exists(title)
	if err != nil {
		ERROR("No todo list entry found with title " + title)
		return
	}

	content = wordwrap.WrapString(content, uint(maxWidth))
	fmt.Printf("%sOld content for %s:\n", defaultColor, title)
	fmt.Printf("%s%s\033[0m\n\n", contentColor, content)

	reader := bufio.NewReader(os.Stdin)
	if conf.Keep_content {
		fmt.Printf("%sEnter content to be added into todo titled %s: \033[0m\n", wipColor, title)
	} else {
		fmt.Printf("%sEnter new content for todo titled %s: \033[0m\n", wipColor, title)
	}

	fmt.Printf("%s❯ \033[0m", borderColor)
	newContent, err := reader.ReadString('\n')
	if err != nil {
		ERROR("Error reading input, did text end with a newline?")
		Usage_edit()
		return
	}
	newContent = strings.TrimSpace(newContent)
	if conf.Keep_content {
		newContent = content + "\n\nNew edit:\n" + newContent
	}
	change_entry_content(newContent, title)
}

// Send the new content of an entry into the database
func change_entry_content(newContent string, title string) {
	todoDB := open_todo_db()
	defer todoDB.Close()

	sqlStatement := `UPDATE todo SET content = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, newContent, title)
	if err != nil {
		ERROR("Failed to edit todo content")
		panic(err)
	}

	OK("Successfully changed content for " + title)
}

// NOTE: Edit command help and usage functions
func Usage_edit() {
	fmt.Print(`Usage: cldl edit [-h | --help] [-k | --keep] <title>
    Use 'cldl edit --help' to see more
`)
}

const HelpEdit = `Help for cldl edit:
    Available arguments:
        --help, -h  | Show this message
        --keep, -k  | Toggles the behavior of keeping on edit

    Edit an existing todo list entry with a given title. Can be configured
    to either keep the existing content or override it by default

    Same text inputting rules apply for editing as adding a new entry,
    Check 'cldl add --help'.`
