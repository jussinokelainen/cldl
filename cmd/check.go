package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Check_todos(confirm_rm bool, check_directories bool, dir_paths []string, verbose bool) {
	if check_directories {
		check_dir_list(dir_paths, verbose)
	} else {
		check_loc_list(confirm_rm)
	}
}

func check_dir_list(dir_paths []string, verbose bool) {
	INFO("Checking for hidden cldl lists")
	if len(dir_paths) < 1 {
		INFO("No paths specified to check, nothing to do.")
		return
	}

	var found []string
	for _, base_dir := range dir_paths {
		// Expand tildes into full file path
		if strings.HasPrefix(base_dir, "~") {
			homedir, err := os.UserHomeDir()
			if err != nil {
				ERROR("Failed to fetch user home dir")
				panic(err)
			}
			base_dir = strings.Replace(base_dir, "~", homedir, 1)
		}

		if !File_exists(base_dir) {
			INFO("This specified base directory doesn't exist or was not found:", base_dir)
			continue
		}
		err := filepath.WalkDir(base_dir, func(path string, dest fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !dest.IsDir() && dest.Name() == ".cldl.db" {
				found = append(found, path)
			}

			return nil
		})

		if err != nil {
			ERROR(err)
			return
		}
	}

	for _, path := range found {
		masterdb_exists := check_if_masterdb_has_loc(path)
		homedir, err := os.UserHomeDir()
		if err != nil {
			ERROR("Failed to fetch user home dir")
			panic(err)
		}
		shortenedString := strings.Replace(path, homedir, "~", 1)
		if masterdb_exists {
			if verbose {
				OK("cldl list at\033[32m", shortenedString, "\033[0mis already in saved locations")
			}
		} else {
			INFO("cldl list at\033[31m", shortenedString, "\033[0mwas not in saved locations, adding...")
			add_to_master_db(path)
		}
	}
	OK("Checks done.")
}

func check_loc_list(confirm_rm bool) {
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
				if ask_yes_no("Do you want to remove it from the list?") {
					remove_from_master_db(location)
				}

			} else {
				remove_from_master_db(location)
			}
		}
	}
	OK("Checking successful")
}

// NOTE: Check command help and usage functions
func Usage_check() {
	fmt.Print(`Usage: cldl check [-h | --help] [--no-confirm]
    Use 'cldl check --help' to see more
`)
}

const HelpCheck = `Help for cldl check:
    Available arguments:
        --help, -h        | Show this message
        --no-confirm      | Don't ask for confirmation before deleting
                          | local databases, regardless of what configs have
        --directories, -d | check directories that are defined in configs for
                          | hidden lists
        --verbose, -v     | make checks more verbose

    Checks all the saved locations of local todo lists and looks for locations
    that do not have the corresponding list files, and asks to delete them`
