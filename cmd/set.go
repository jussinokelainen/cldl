package cmd

import "fmt"

func Set_filepath_to_entry(title string, file_path string) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	_, exists := get_content_if_entry_exists(title)
	if exists != nil {
		ERROR(fmt.Sprintf("No todo entry exists with title %s", title))
		return
	}

	INFO("Setting filepath of", title, "to", file_path)
	sqlStatement := `UPDATE todo SET file = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, file_path, title)
	if err != nil {
		ERROR("Failed to edit todo content")
		panic(err)
	}
}
func Set_fileline_to_entry(title string, line_num int) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	_, exists := get_content_if_entry_exists(title)
	if exists != nil {
		ERROR(fmt.Sprintf("No todo entry exists with title %s", title))
		return
	}

	INFO("Setting line number of", title, "to", line_num)
	sqlStatement := `UPDATE todo SET line = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, line_num, title)
	if err != nil {
		ERROR("Failed to edit todo content")
		panic(err)
	}
}

func EditPriority(title string, newPrio int) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	_, exists := get_content_if_entry_exists(title)
	if exists != nil {
		ERROR(fmt.Sprintf("No todo entry exists with title %s", title))
		return
	}

	INFO(fmt.Sprintf("Setting priority of %s to %d", title, newPrio))
	sqlStatement := `UPDATE todo SET priority = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, newPrio, title)
	if err != nil {
		ERROR("Failed to edit todo content")
		panic(err)
	}
}

func SetTagToEntry(title string, tag string) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	_, exists := get_content_if_entry_exists(title)
	if exists != nil {
		ERROR(fmt.Sprintf("No todo entry exists with title %s", title))
		return
	}

	INFO(fmt.Sprintf("Setting tag of %s to %s", title, tag))
	sqlStatement := `UPDATE todo SET tag = $1 WHERE UPPER(title) = UPPER($2);`
	_, err := todoDB.Exec(sqlStatement, tag, title)
	if err != nil {
		ERROR("Failed to edit todo content")
		panic(err)
	}
}

func UsageSet() {
	fmt.Print(`Usage: todo set [-h | --help] [-p | --priority] <title>
    Use 'todo set --help' to see more
`)
}

const HelpSet = `Help for todo set:
    Available arguments:
        --help, -h     | Show this message
        --priority, -p | Set the priority of an entry
        --tag, -t      | Set the tag of an entry
        --file, -f     | Set the file of an entry
        --line, -l     | Set the file's line number of an entry

    Set various things in already existing entries should be simple`
