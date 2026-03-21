package cmd

import (
	"fmt"
	"os"
)

func RmTodo(title string, rmAll bool, conf RmConf) {
	if rmAll {
		if getEntryCount() != 0 {
			info("The list is not empty!")
			if askYesNo("Do you still want to remove it?") {
				removeAllData()
			}
		} else if conf.Always_confirm_full {
			if askYesNo("Are you sure?") {
				removeAllData()
			}
		} else {
			removeAllData()
		}
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
		if conf.Ask_full && getEntryCount() == 0 {
			info("The last entry of this todo-list was removed.")
			if askYesNo("Do you want to fully remove the list?") {
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

// Deletes the local todo database and removes it from the master list
// If there is a local file with the exact name that is not a todo database:
// don't care + didn't ask + skill issue + your file is deleted
func removeAllData() {
	todoPath := GetDbPath()
	removeFromMasterDB(todoPath)
	err := os.Remove(todoPath)
	if err != nil {
		errout("Failed removing local db")
		panic(err)
	}

	ok("Succesfully removed local database!")
}

// NOTE: Remove command help and usage functions
func UsageRm() {
	fmt.Print(`Usage: todo rm [-h | --help] [-a | --all] <title>
    Use 'todo rm --help' to see more
`)
}
func HelpRm() {
	const helpmsg = `Help for todo rm / done:
    Available arguments:
        --help, -h  | Show help for todo rm
        --all, -a   | Fully remove todo list from current directory

    Rm and done are the same command with a different name.
    Use 'todo rm <title>' where <title> is the title
    for the list entry to be deleted.`

	PrintHelpMSG(helpmsg)
}
