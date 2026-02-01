package main

import (
	"os"

	_ "modernc.org/sqlite"
)

func initTodo(args []string) {
	if len(args) < 1 {
		createNewTodo()
		return
	}

	switch args[0] {
	case "--help", "-h":
		helpInit()
		return
	default:
		errout("Invalid arguments")
		usageInit()
	}
}

// Create a new todo into cwd if one doesn't exist already
func createNewTodo() {
	todoPath := getDbPath()
	// Check whether current directory already has a list
	// if it exists, do not create a new one, and just return
	if _, err := os.Stat(getDbPath()); !os.IsNotExist(err) {
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
	_, err = masterDB.Exec(sqlStatement, todoPath)
	if err != nil {
		errout("Adding to master DB failed!")
		panic(err)
	}

	ok("New todo created!")
}
