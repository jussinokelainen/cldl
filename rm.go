package main

import "os"

func rmTodo(args []string) {
	if !todoExists() {
		return
	}

	if len(args) < 1 {
		usageRm()
		return
	}

	switch args[0] {
	case "--help", "-h":
		helpRm()
		return
	case "--all", "-a":
		removeAllData()
		return
	default:
		title := args[0]
		todoDB := openTodoDB()
		defer todoDB.Close()

		sqlStatement := `DELETE FROM todo WHERE title= ?;`
		res, err := todoDB.Exec(sqlStatement, title)
		if err != nil {
			errout("Error removing entry")
			panic(err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			errout("Error getting affected rows")
			panic(err)
		}

		if rowsAffected == 0 {
			errout("No entry found with title " + title)
		} else {
			ok("Succesfully removed entry " + title)
		}
	}
}

// Deletes the local todo database and removes it from the master list
// If there is a local file with the exact name that is not a todo database:
// don't care + didn't ask + skill issue + your file is deleted
func removeAllData() {
	todoPath := getDbPath()
	sqlStatement := `DELETE FROM locations WHERE location = ?;`
	_, err := masterDB.Exec(sqlStatement, todoPath)
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
