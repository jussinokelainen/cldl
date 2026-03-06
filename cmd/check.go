package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func CheckTodos(confirm_rm bool) {
	var locSlice []string

	rows, err := MasterDB.Query(`SELECT * FROM locations;`)
	if err != nil {
		errout("Failed getting all locations")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var location string
		err = rows.Scan(&location)
		if err != nil {
			errout("Failed scanning entry content")
			panic(err)
		}

		locSlice = append(locSlice, location)
	}

	for _, location := range locSlice {
		if _, err := os.Stat(location); os.IsNotExist(err) {
			info("This todo does not exist: \n\033[35m  " + location + "\033[0m")
			if confirm_rm {
				if confirmMasterRm() {
					remove_master_entry(location)
				}

			} else {
				remove_master_entry(location)
			}
		}
	}
}

func confirmMasterRm() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to remove it from the list? [y/n]: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		errout("Error reading input")
		return confirmMasterRm()
	}
	answer = strings.TrimSpace(answer)
	switch answer {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Print("Invalid answer, try again.\n")
		return confirmMasterRm()
	}
}

// NOTE: Check command help and usage functions
func UsageCheck() {
	fmt.Print(`
Usage: todo check
    Use 'todo check -help' to see more
`)
}
func HelpCheck() {
	fmt.Print(`
Help for todo check:
    Available arguments:
        --help, -h   | Show help for todo check
        --no-confirm | Don't ask for confirmation before deleting
                       local databases, regardless of what configs have

    Checks all the saved locations of local todo lists and looks
    for locations that do not have the corresponding list files.

    Config option 'ask_rm_on_check' determines whether excess locations
    are automatically removed or prompted
`)
}
