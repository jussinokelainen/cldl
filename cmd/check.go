package cmd

import (
	"fmt"
	"os"
)

func CheckTodos(confirm_rm bool) {
	var locSlice []string

	rows, err := MasterDB.Query(`SELECT location FROM locations;`)
	if err != nil {
		ERROR("Failed getting all locations")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var location string
		err = rows.Scan(&location)
		if err != nil {
			ERROR("Failed scanning entry content")
			panic(err)
		}

		locSlice = append(locSlice, location)
	}

	for _, location := range locSlice {
		if _, err := os.Stat(location); os.IsNotExist(err) {
			INFO("This todo does not exist: \n\033[35m  " + location + "\033[0m")
			if confirm_rm {
				if askYesNo("Do you want to remove it from the list?") {
					removeFromMasterDB(location)
				}

			} else {
				removeFromMasterDB(location)
			}
		}
	}
	OK("Checking successful")
}

// NOTE: Check command help and usage functions
func UsageCheck() {
	fmt.Print(`Usage: todo check [-h | --help] [--no-confirm]
    Use 'todo check --help' to see more
`)
}

const HelpCheck = `Help for todo check:
    Available arguments:
        --help, -h   | Show this message
        --no-confirm | Don't ask for confirmation before deleting
                       local databases, regardless of what configs have

    Checks all the saved locations of local todo lists and looks for locations
    that do not have the corresponding list files, and asks to delete them`
