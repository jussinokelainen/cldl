package cmd

import (
	"fmt"
	"os"
)

func RmTodo(title string, rmAll bool) {
	if rmAll {
		removeAllData()
		return
	}

	if title == "" {
		UsageRm()
		return
	}

	todoDB := openTodoDB()
	defer todoDB.Close()

	sqlStatement := `DELETE FROM todo WHERE UPPER(title) = UPPER(?);`
	res, err := todoDB.Exec(sqlStatement, title)
	if err != nil {
		errout("Error removing entry")
		panic(err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		errout("No entry found with title " + title)
	} else {
		ok("Succesfully removed entry " + title)
	}
}

// Deletes the local todo database and removes it from the master list
// If there is a local file with the exact name that is not a todo database:
// don't care + didn't ask + skill issue + your file is deleted
func removeAllData() {
	todoPath := GetDbPath()
	sqlStatement := `DELETE FROM locations WHERE location = ?;`
	_, err := MasterDB.Exec(sqlStatement, todoPath)
	if err != nil {
		errout("Error removing from master db")
		panic(err)
	}
	err = os.Remove(todoPath)
	if err != nil {
		errout("Failed removing local db")
		panic(err)
	}

	ok("Succesfully removed database!")
}

// NOTE: Remove command help and usage functions
func UsageRm() {
	fmt.Print(`
Usage: todo rm [<args>] <title>
	Use 'todo rm --help' to see arguments
`)
}
func HelpRm() {
	fmt.Print(`
Help for todo rm / done:
	Available arguments:
		-help  | Show help for todo rm
		-h     | Same as '--help'
		-all   | Fully remove todo list from current directory
		-a     | Same as '--all'

	Rm and done are the same command with a different name.
	Use 'todo rm <title>' where <title> is the title
	for the list entry to be deleted.
`)
}
