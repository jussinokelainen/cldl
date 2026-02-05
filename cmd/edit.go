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
		errout("Error reading input")
		panic(err)
	}
	newContent = strings.TrimSpace(newContent)
	if keep {
		newContent = content + "\n\nNew edit:\n" + newContent
	}
	changeEntryContent(newContent, title)
}

// Get the existing content of a todo entry if it exists, returns an error if it doesn't exist
func getIfEntryExists(title string) (string, error) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	sqlStatement := `SELECT content from todo WHERE title = $1;`
	res, err := todoDB.Query(sqlStatement, title)
	if err != nil {
		errout("Failed checking entry!")
		panic(err)
	}
	defer res.Close()

	var content string
	for res.Next() {
		err = res.Scan(&content)
		if err != nil {
			errout("Row scanning failed")
			panic(err)
		}
	}
	if content == "" {
		return content, fmt.Errorf("No content found")
	}

	return content, nil
}

// Send the new content of an entry into the database
func changeEntryContent(newContent string, title string) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	sqlStatement := `UPDATE todo SET content = $1 WHERE title = $2`
	_, err := todoDB.Exec(sqlStatement, newContent, title)
	if err != nil {
		errout("Failed to edit todo content")
		panic(err)
	}

	fmt.Print("\n")
	ok("Successfully changed content for " + title)

}
