package cmd

import (
	"database/sql"
	"fmt"
	"os"
)

var MasterDB *sql.DB

// Creates a 'master' database that holds all the locations to
// local databases if one does not yet exist
func CreateMasterDB() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		errout("Getting homedir failed!")
		panic(err)
	}
	masterDbDIR := homedir + "/.sqlite/todo"
	err = os.MkdirAll(masterDbDIR, 0755)
	if err != nil {
		errout("Creating .sqlite/todo directory failed!")
		panic(err)
	}

	MasterDB, err = sql.Open("sqlite", masterDbDIR+"/.todo.db")
	if err != nil {
		errout("Opening master database failed!")
		panic(err)
	}
	_, err = MasterDB.Exec(`CREATE TABLE IF NOT EXISTS locations (location VARCHAR UNIQUE);`)
	if err != nil {
		errout("Creating master database failed!")
		panic(err)
	}
}

// Get the path of a local todo database, returns the path as a string
func GetDbPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		errout("Getting db path failed!")
		panic(err)
	}
	return cwd + "/.todoApp.db"
}

// Open a connection to a local database, returns a pointer to it
func openTodoDB() *sql.DB {
	db, err := sql.Open("sqlite", GetDbPath())
	if err != nil {
		errout("Opening todo DB failed!")
		panic(err)
	}
	return db
}

// Checks whether a local todo exists in the current directory,
// might give an erroneus result if some other file is named
// exactly as the todo database should be,
// in which case errors that might come are a skill issue
func TodoExists() bool {
	if _, err := os.Stat(GetDbPath()); os.IsNotExist(err) {
		return false
	}
	return true
}

// Status printing helpers
func ok(msg string)     { fmt.Println("[\033[32m OK \033[0m] ", msg) }
func info(msg string)   { fmt.Println("[\033[35m INFO \033[0m] ", msg) }
func errout(msg string) { fmt.Println("[\033[31m ERROR \033[0m] ", msg) }
