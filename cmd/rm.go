package cmd

import (
	"fmt"
	"os"
)

func Rm_todo(title string, rmAll bool, rmTag bool, conf RmConf) {
	if rmAll && rmTag {
		todoDB := open_todo_db()
		defer todoDB.Close()

		INFO("Clearing all tags in current todo list")
		_, err := todoDB.Exec("UPDATE todo SET tag = 'NONE';")
		if err != nil {
			ERROR("Failed to edit todo content")
			panic(err)
		}
		return
	} else if rmTag {
		if title == "" {
			Usage_rm()
			return
		}
		Set_tag_to_entry(title, "NONE")

		return
	} else if rmAll {
		if get_entry_count() != 0 {
			INFO("The list is not empty!")
			if ask_yes_no("Do you still want to remove it?") {
				remove_all_data()
			}
		} else if conf.Always_confirm_full {
			if ask_yes_no("Are you sure?") {
				remove_all_data()
			}
		} else {
			remove_all_data()
		}
		return
	}

	if title == "" {
		Usage_rm()
		return
	}

	todoDB := open_todo_db()
	defer todoDB.Close()

	sqlStatement := `DELETE FROM todo WHERE UPPER(title) = UPPER(?);`
	res, err := todoDB.Exec(sqlStatement, title)
	if err != nil {
		ERROR("Error removing entry")
		panic(err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		ERROR("No entry found with title " + title)
	} else {
		OK("Succesfully removed entry " + title)
		if conf.Ask_full && get_entry_count() == 0 {
			INFO("The last entry of this todo-list was removed.")
			if ask_yes_no("Do you want to fully remove the list?") {
				remove_all_data()
			}
		}
	}
}

func get_entry_count() int {
	var count int
	todoDB := open_todo_db()
	defer todoDB.Close()

	res, err := todoDB.Query(`SELECT COUNT(*) FROM todo;`)
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&count)
		if err != nil {
			ERROR("Failed scanning existing entry content")
			panic(err)
		}
	}

	return count
}

// Deletes the local todo database and removes it from the master list
// If there is a local file with the exact name that is not a todo database:
// don't care + didn't ask + skill issue + your file is deleted
func remove_all_data() {
	todoPath := Get_db_path()
	remove_from_master_db(todoPath)
	err := os.Remove(todoPath)
	if err != nil {
		ERROR("Failed removing local db")
		panic(err)
	}

	OK("Succesfully removed local database!")
}

// NOTE: Remove command help and usage functions
func Usage_rm() {
	fmt.Print(`Usage: cldl rm [-h | --help] [-a | --all] <title>
    Use 'cldl rm --help' to see more
`)
}

const HelpRm = `Help for cldl rm / done:
    Available arguments:
        --help, -h  | Show help for cldl rm
        --all, -a   | Fully remove cldl list from current directory
        --tag, -t   | Remove the tag from an entry. If used together with
                    | --all, clears all tags in current list
        --file, -f  | Clear the file set for an entry

    Rm and done are the same command with a different name.
    Use 'cldl rm <title>' where <title> is the title
    for the list entry to be deleted.`
