package cmd

import (
	"fmt"
)

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

func UsagePriority() {
	fmt.Print(`
Usage: todo set-priority <number> <title>
	Use 'todo priority --help' to see more
`)
}
func HelpPriority() {
	fmt.Print(`
Help for todo priority:
	Available arguments:
		--help, -h   | Show this message

    Set the priority of a todo entry
`)

}
