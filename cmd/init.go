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
		ERROR("Todo already exists in current directory!")
		return
	}

	todoDB := openTodoDB()
	defer todoDB.Close()
	_, err := todoDB.Exec(`CREATE TABLE todo(
        title VARCHAR UNIQUE NOT NULL,
        content VARCHAR NOT NULL,
        time INTEGER NOT NULL,
        priority INTEGER NOT NULL,
        tag VARCHAR NOT NULL,
		file VARCHAR NOT NULL,
		line INTEGER NOT NULL
        );`)
	if err != nil {
		ERROR("Creating new todo failed!")
		panic(err)
	}

	addToMasterDB(todoPath)

	OK("New todo created!")
}

// NOTE: Init help and usage functions
func UsageInit() {
	fmt.Print(`Default usage: cldl init [-h | --help]
    Use 'cldl init --help' to see more
`)
}

const HelpInit = `Help for cldl init:
    Available arguments:
        --help, -h  | Show this message

    Initialize a local todo list in current directory
    Might be a useless command, since initialization can also
    be done when adding an entry`
