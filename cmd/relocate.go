package cmd

import "fmt"

func RelocateTodo(ask_rm_on_check bool) {
	path := GetDbPath()
	sqlStatement := `SELECT COUNT(*) from locations WHERE location = ?;`
	res, err := MasterDB.Query(sqlStatement, path)
	if err != nil {
		ERROR("Failed checking entry in database!")
		panic(err)
	}
	defer res.Close()

	var content int
	for res.Next() {
		err = res.Scan(&content)
		if err != nil {
			ERROR("Failed scanning existing entry content")
			panic(err)
		}
	}
	if content < 1 {
		INFO("Adding local todo to location list")
		addToMasterDB(path)
	} else {
		INFO("Todo location already exists in the list")
	}
	CheckTodos(ask_rm_on_check)
}

// NOTE: Init help and usage functions
func UsageRelocate() {
	fmt.Print(`Default usage: cldl relocate [-h | --help]
    Use 'cldl relocate --help' to see more
`)
}

const HelpRelocate = `Help for cldl relocate:
    Available arguments:
        --help, -h  | Show this message

    If a cldl list exists in current location but
    it isn't present in the location list, adds
    it and checks the location list for non-existent
    locations.

    Useful when renaming directories etc.`
