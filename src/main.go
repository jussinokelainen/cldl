package main

import (
	"database/sql"
	"fmt"
	"os"
)

var masterDB *sql.DB

func main() {
	args := os.Args[1:]
	// If no args given, print usage and exit
	if len(args) < 1 {
		mainUsage()
		return
	}

	createMasterDB()
	defer masterDB.Close()

	switch args[0] {
	case "--help", "-h":
		mainHelp()
	case "init":
		initTodo(args[1:])
	case "add":
		if !todoExists() {
			return
		}
		addTodo(args[1:])
	case "list":
		listTodo(args[1:])
	case "rm", "remove", "done":
		if !todoExists() {
			return
		}
		rmTodo(args[1:])
	case "edit":
		if !todoExists() {
			return
		}
		editTodo(args[1:])
	default:
		errout("Bad arguments")
		mainUsage()
	}
}

// Status printing helpers
func ok(msg string)     { fmt.Println("[\033[32m OK \033[0m] ", msg) }
func info(msg string)   { fmt.Println("[\033[35m INFO \033[0m] ", msg) }
func errout(msg string) { fmt.Println("[\033[31m ERROR \033[0m] ", msg) }

// Get the path of a local todo database, returns the path as a string
func getDbPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		errout("Getting db path failed!")
		panic(err)
	}
	return cwd + "/.todoApp.db"
}

// Open a connection to a local database, returns a pointer to it
func openTodoDB() *sql.DB {
	db, err := sql.Open("sqlite", getDbPath())
	if err != nil {
		errout("Opening todo DB failed!")
		panic(err)
	}
	return db
}

// Creates a 'master' database that holds all the locations to
// local databases if one does not yet exist
func createMasterDB() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		errout("Getting homedir failed!")
		panic(err)
	}
	masterDbDIR := homedir + "/.sqlite/todo"
	err = os.MkdirAll(masterDbDIR, 0755)
	if err != nil {
		errout("Creating .sqlite/todo dir failed!")
		panic(err)
	}

	masterDB, err = sql.Open("sqlite", masterDbDIR+"/.todo.db")
	if err != nil {
		errout("Opening master db failed!")
		panic(err)
	}
	_, err = masterDB.Exec(`CREATE TABLE IF NOT EXISTS locations (location VARCHAR UNIQUE);`)
	if err != nil {
		errout("Creating master db failed!")
		panic(err)
	}
}

// Checks whether a local todo exists in the current directory,
// might give an erroneus result if some other file is named
// exactly as the todo database should be,
// in which case errors that might come are a skill issue
func todoExists() bool {
	if _, err := os.Stat(getDbPath()); os.IsNotExist(err) {
		errout("No todo exists in current directory!")
		return false
	}
	return true
}
