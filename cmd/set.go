package cmd

import "fmt"

func EditPriority(title string, newPrio int) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	_, exists := getIfEntryExists(title)
	if exists != nil {
		errout(fmt.Sprintf("No todo entry exists with title %s", title))
		return
	}

	info(fmt.Sprintf("Setting priority of %s to %d", title, newPrio))
	sqlStatement := `UPDATE todo SET priority = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, newPrio, title)
	if err != nil {
		errout("Failed to edit todo content")
		panic(err)
	}
}

func SetTagToEntry(title string, tag string) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	_, exists := getIfEntryExists(title)
	if exists != nil {
		errout(fmt.Sprintf("No todo entry exists with title %s", title))
		return
	}

	info(fmt.Sprintf("Setting tag of %s to %s", title, tag))
	sqlStatement := `UPDATE todo SET tag = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, tag, title)
	if err != nil {
		errout("Failed to edit todo content")
		panic(err)
	}
}

func UsageSet() {
	fmt.Print(`
Usage: todo set [-h | --help] [-p | --priority] <title>
    Use 'todo set --help' to see more
`)
}
func HelpSet() {
	fmt.Print(`
Help for todo set:
    Available arguments:
        --help, -h     | Show this message
        --priority, -p | Set the priority of an entry
        --tag, -t      | Set the tag of an entry

    Set various things in already existing entries
`)

}
