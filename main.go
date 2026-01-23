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
	default:
		errout("Bad arguments")
		mainUsage()
	}
}

func ok(msg string)     { fmt.Println("[\033[32m OK \033[0m] ", msg) }
func info(msg string)   { fmt.Println("[\033[35m INFO \033[0m] ", msg) }
func errout(msg string) { fmt.Println("[\033[31m ERROR \033[0m] ", msg) }

func getDbPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		errout("Getting db path failed!")
		panic(err)
	}
	return cwd + "/.todoApp.db"
}

func openTodoDB() *sql.DB {
	db, err := sql.Open("sqlite", getDbPath())
	if err != nil {
		errout("Opening todo DB failed!")
		panic(err)
	}
	return db
}

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

func todoExists() bool {
	if _, err := os.Stat(getDbPath()); os.IsNotExist(err) {
		errout("No todo exists in current directory!")
		return false
	}
	return true
}
