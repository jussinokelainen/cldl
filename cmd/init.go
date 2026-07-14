package cmd

import (
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func Init_todo() {
	todoPath := Get_db_path()
	// Check whether current directory already has a list
	// if it exists, do not create a new one, and just return
	if _, err := os.Stat(Get_db_path()); !os.IsNotExist(err) {
		ERROR("Todo already exists in current directory!")
		return
	}

	todoDB := open_todo_db()
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

	add_to_master_db(todoPath)

	OK("New todo created!")
}

// NOTE: Init help and usage functions
func Usage_init() {
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
