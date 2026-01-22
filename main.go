package main

import (
	"database/sql"
	"fmt"
	"os"
)

var masterDB *sql.DB

func getDbPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		error("Getting db path failed!")
		panic(err)
	}
	return cwd + "/.todoApp.db"
}

func info(msg string) {
	fmt.Println("[\033[35m INFO \033[0m] ", msg)
}

func error(msg string) {
	fmt.Println("[\033[31m ERROR \033[0m] ", msg)
}

func printUsage() {
	fmt.Println("Usage: todo <COMMAND> [<args>]")
	fmt.Println("  Use todo --help to see arguments")
}

func printHelp() {
	fmt.Println("Help for todo:")
	fmt.Println(" Usage: todo <COMMAND> [<args>]")
	fmt.Println(" Available commands:")
	fmt.Println("   init | create new todo in current dir")
}

func openTodoDB() *sql.DB {
	db, err := sql.Open("sqlite", getDbPath())
	if err != nil {
		error("Opening todo DB failed!")
		panic(err)
	}
	return db
}

func createMasterDB() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		error("Getting homedir failed!")
		panic(err)
	}

	masterDbDIR := homedir + "/.sqlite/todo"
	err = os.MkdirAll(masterDbDIR, 0755)
	if err != nil {
		error("Creating .sqlite/todo dir failed!")
		panic(err)
	}

	masterDB, err = sql.Open("sqlite", masterDbDIR+"/.todo.db")
	if err != nil {
		error("Opening master db failed!")
		panic(err)
	}

	_, err = masterDB.Exec(`CREATE TABLE IF NOT EXISTS locations (location VARCHAR UNIQUE);`)
	if err != nil {
		error("Creating master db failed!")
		panic(err)
	}
}

func main() {
	args := os.Args[1:]
	createMasterDB()
	defer masterDB.Close()

	// If no args given, print usage and exit
	if len(args) < 1 {
		printUsage()
		return
	}

	switch args[0] {
	case "--help":
		printHelp()

	case "init":
		initTodo(args[1:])

	case "add":
		addTodo(args[1:])

	case "list":
		listTodo(args[1:])

	case "rm":
		rmTodo(args[1:])

	default:
		info("Bad arguments")
		printUsage()
	}
}
