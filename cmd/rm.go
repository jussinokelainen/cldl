package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RmTodo(title string, rmAll bool, ask_rm_all bool) {
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
		if ask_rm_all && getEntryCount() == 0 {
			info("The last entry of this todo-list was removed.")
			if ask_full_rm() {
				removeAllData()
			}
		}
	}
}

func getEntryCount() int {
	var count int
	todoDB := openTodoDB()
	defer todoDB.Close()

	res, err := todoDB.Query(`SELECT COUNT(*) FROM todo;`)
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&count)
		if err != nil {
			errout("Failed scanning existing entry content")
			panic(err)
		}
	}

	return count
}

func ask_full_rm() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to fully remove the list? [y/n]: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		errout("Error reading input")
		return ask_full_rm()
	}
	answer = strings.TrimSpace(answer)
	switch answer {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Print("Invalid answer, try again.\n")
		return ask_full_rm()
	}
}

// Deletes the local todo database and removes it from the master list
// If there is a local file with the exact name that is not a todo database:
// don't care + didn't ask + skill issue + your file is deleted
func removeAllData() {
	todoPath := GetDbPath()
	remove_master_entry(todoPath)
	err := os.Remove(todoPath)
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
		--help, -h  | Show help for todo rm
		--all, -a   | Fully remove todo list from current directory

	Rm and done are the same command with a different name.
	Use 'todo rm <title>' where <title> is the title
	for the list entry to be deleted.

    Config option 'ask_full_rm' can determine whether removing the local
    database on the removal of the last entry in it gets asked or not
`)
}
