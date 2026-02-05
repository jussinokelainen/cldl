package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Adds a new todo with a title given as an argument, 'todo add <title>'
func AddTodo(title string) {

	todoDB := openTodoDB()
	defer todoDB.Close()

	reader := bufio.NewReader(os.Stdin)

	time := time.Now().Unix()

	fmt.Printf("\033[36mEnter contents for new todo titled %s: \033[0m\n", title)
	fmt.Print("\033[35m❯ \033[0m")
	content, err := reader.ReadString('\n')
	if err != nil {
		errout("Error reading input")
		panic(err)
	}
	content = strings.TrimSpace(content)

	sqlStatement := `INSERT INTO todo(title, content, time) VALUES($1, $2, $3);`
	_, err = todoDB.Exec(sqlStatement, title, content, time)
	if err != nil {
		errout("Error adding new todo, is a list initialized?")
		panic(err)
	}
	ok("Successfully added new todo " + title)
}
