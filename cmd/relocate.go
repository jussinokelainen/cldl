package cmd

import "fmt"

func RelocateTodo(ask_rm_on_check bool) {
	path := GetDbPath()
	sqlStatement := `SELECT COUNT(*) from locations WHERE location = ?;`
	res, err := MasterDB.Query(sqlStatement, path)
	if err != nil {
		errout("Failed checking entry in database!")
		panic(err)
	}
	defer res.Close()

	var content int
	for res.Next() {
		err = res.Scan(&content)
		if err != nil {
			errout("Failed scanning existing entry content")
			panic(err)
		}
	}
	if content < 1 {
		info("Adding local todo to location list")
		addToMasterDB(path)
	} else {
		info("Todo location already exists in the list")
	}
	CheckTodos(ask_rm_on_check)
}

// NOTE: Init help and usage functions
func UsageRelocate() {
	fmt.Print(`Default usage: todo relocate [-h | --help]
    Use 'todo relocate --help' to see more
`)
}
func HelpRelocate() {
	const helpmsg = `Help for todo relocate:
    Available arguments:
        --help, -h  | Show this message

    If a todo list exists in current location but
    it isn't present in the location list, adds
    it and checks the location list for non-existent
    locations.

    Useful when renaming directories etc.`

	PrintHelpMSG(helpmsg)
}
