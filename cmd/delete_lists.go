package cmd

import (
	"os"
)

func Delete_saved_lists() {
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
			ERROR("Error scanning entry content")
			panic(err)
		}

		locSlice = append(locSlice, location)
	}
	if rows.Err() != nil {
		ERROR("Error scanning entry content")
		panic(err)
	}

	if len(locSlice) < 1 {
		INFO("No list locations saved, nothing to do.")
		os.Exit(0)
	}

	if ask_yes_no("\033[31mRemoving lists CANNOT BE UNDONE.\033[0m Are you sure?") {
		for _, list := range locSlice {
			INFO("Removing", list)
			os.Remove(list)
		}
		homedir, err := os.UserHomeDir()
		if err != nil {
			ERROR("Error getting user home directory")
			os.Exit(1)
		}

		masterDbDIR := homedir + "/.local/share/cldl"
		os.RemoveAll(masterDbDIR)
	}
}
