package cmd

import (
	"database/sql"
	"fmt"
)

func FixTodoTable(defaultPriority int) {
	title := false
	content := false
	time := false
	priority := false
	info("Checking todo table")

	todoDB := openTodoDB()
	defer todoDB.Close()

	rows, err := todoDB.Query("PRAGMA table_info(todo)")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var (
		cid       int
		name      string
		ctype     string
		notnull   int
		dfltValue sql.NullString
		pk        int
	)

	for rows.Next() {
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			panic(err)
		}
		switch name {
		case "title":
			title = true
		case "content":
			content = true
		case "time":
			time = true
		case "priority":
			priority = true
		default:
			errout("Unexpected column name found: " + name)
		}
	}
	if !title {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN title VARCHAR UNIQUE NOT NULL;`)
		if err != nil {
			panic(err)
		}
		info("Added missing title column")
	}
	if !content {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN content VARCHAR NOT NULL;`)
		if err != nil {
			panic(err)
		}
		info("Added missing content column")
	}
	if !time {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN time INTEGER NOT NULL;`)
		if err != nil {
			panic(err)
		}
		info("Added missing time column")
	}
	if !priority {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN priority INTEGER NOT NULL;`)
		if err != nil {
			panic(err)
		}
		info("Added missing priority column")
	}
	DefaultNullPriorities(defaultPriority)
}

func DefaultNullPriorities(defaultPriority int) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	info("Setting possible null priority values to default value")
	sqlStatement := `UPDATE todo SET priority = $1 WHERE priority IS NULL;`
	_, err := todoDB.Exec(sqlStatement, defaultPriority)
	if err != nil {
		errout("Failed to edit todo content")
		panic(err)
	}
}

func UsageFix() {
	fmt.Print(`
Usage: todo fix
	Use 'todo fix --help' to see more
`)
}
func HelpFix() {
	fmt.Print(`
Help for todo fix:
	Available arguments:
		--help, -h   | Show this message

    This command will be useful after making breaking
    changes to the program and its databases. It whether
    a local todo list has all the required columns, and
    adds them if they don't exist.
`)

}
