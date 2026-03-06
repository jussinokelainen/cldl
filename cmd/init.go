package cmd

import (
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func InitTodo() {
	todoPath := GetDbPath()
	// Check whether current directory already has a list
	// if it exists, do not create a new one, and just return
	if _, err := os.Stat(GetDbPath()); !os.IsNotExist(err) {
		errout("Todo already exists in current directory!")
		return
	}

	todoDB := openTodoDB()
	defer todoDB.Close()
	_, err := todoDB.Exec(`CREATE TABLE todo(
		title VARCHAR UNIQUE,
		content VARCHAR,
		time INTEGER
		);`)
	if err != nil {
		errout("Creating new todo failed!")
		panic(err)
	}

	// Add new todo location into list location database
	sqlStatement := `INSERT INTO locations(location) VALUES($1);`
	_, err = MasterDB.Exec(sqlStatement, todoPath)
	if err != nil {
		errout("Adding to master DB failed!")
		panic(err)
	}

	ok("New todo created!")
}

// NOTE: Init help and usage functions
func UsageInit() {
	fmt.Print(`
Default usage: todo init
	Use 'todo init -help' to see more
`)
}
func HelpInit() {
	fmt.Print(`
Help for todo init:
	Available arguments:
		--help  | Show help for todo init
		-h      | Same as '--help'

	Initialize a local todo list in current directory
`)
}
