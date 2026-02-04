package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
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
		initFlags := flag.NewFlagSet("initFlags", flag.ExitOnError)
		help := initFlags.Bool("h", false, "show help for todo init")
		helpLong := initFlags.Bool("help", false, "show help for todo init")
		initFlags.Usage = usageInit
		err := initFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			helpInit()
			return
		}
		if len(initFlags.Args()) > 0 {
			fmt.Println("Bad Arguments")
			initFlags.Usage()
			os.Exit(1)
		}
		initTodo()

	case "add":
		addFlags := flag.NewFlagSet("addFlags", flag.ExitOnError)
		help := addFlags.Bool("h", false, "show help for todo add")
		helpLong := addFlags.Bool("help", false, "show help for todo add")
		addFlags.Usage = usageAdd
		err := addFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			helpAdd()
			return
		}
		if !todoExists() {
			return
		}
		title := strings.Join(addFlags.Args(), " ")
		addTodo(title)

	case "list":
		listFlags := flag.NewFlagSet("listFlags", flag.ExitOnError)
		help := listFlags.Bool("h", false, "show help for todo list")
		helpLong := listFlags.Bool("help", false, "show help for todo list")
		pager := listFlags.Bool("pager", false, "Do not send local list to pager")
		all := listFlags.Bool("a", false, "List all todo list locations")
		allLong := listFlags.Bool("all", false, "List all todo list locations")
		listFlags.Usage = usageList
		err := listFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			helpList()
			return
		}
		listAll := false
		if *all || *allLong {
			listAll = true
		}
		// Reverse value since bool flags don't automatically toggle to false
		listTodo(listAll, !*pager)

	case "rm", "remove", "done":
		if !todoExists() {
			return
		}
		rmFlags := flag.NewFlagSet("rmFlags", flag.ExitOnError)
		help := rmFlags.Bool("h", false, "show help for todo rm")
		helpLong := rmFlags.Bool("help", false, "show help for todo rm")
		all := rmFlags.Bool("a", false, "List all todo list locations")
		allLong := rmFlags.Bool("all", false, "List all todo list locations")
		rmFlags.Usage = usageRm
		err := rmFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			helpRm()
			return
		}
		rmAll := false
		if *all || *allLong {
			rmAll = true
		}
		title := strings.Join(rmFlags.Args(), " ")
		rmTodo(title, rmAll)

	case "edit":
		if !todoExists() {
			return
		}
		editFlags := flag.NewFlagSet("editFlags", flag.ExitOnError)
		help := editFlags.Bool("h", false, "show help for todo edit")
		helpLong := editFlags.Bool("help", false, "show help for todo edit")
		keep := editFlags.Bool("k", false, "List all todo list locations")
		keepLong := editFlags.Bool("keep", false, "List all todo list locations")
		editFlags.Usage = usageEdit
		err := editFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			helpEdit()
			return
		}
		keepContent := false
		if *keep || *keepLong {
			keepContent = true
		}
		title := strings.Join(editFlags.Args(), " ")
		editTodo(title, keepContent)

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
