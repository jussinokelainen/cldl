package cmd

import (
	"database/sql"
	"fmt"
)

func Fix_todo_table(defaultPriority int) {
	title := false
	content := false
	time := false
	priority := false
	tag := false
	file := false
	line := false
	INFO("Checking todo table")

	todoDB := open_todo_db()
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
		case "tag":
			tag = true
		case "file":
			tag = true
		case "line":
			tag = true
		default:
			ERROR("Unexpected column name found: " + name)
		}
	}
	if !title {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN title VARCHAR UNIQUE NOT NULL;`)
		if err != nil {
			panic(err)
		}
		INFO("Added missing title column")
	}
	if !content {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN content VARCHAR NOT NULL;`)
		if err != nil {
			panic(err)
		}
		INFO("Added missing content column")
	}
	if !time {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN time INTEGER NOT NULL;`)
		if err != nil {
			panic(err)
		}
		INFO("Added missing time column")
	}
	if !priority {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN priority INTEGER NOT NULL;`)
		if err != nil {
			panic(err)
		}
		INFO("Added missing priority column")
	}
	if !tag {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN tag VARCHAR NOT NULL DEFAULT "NONE";`)
		if err != nil {
			panic(err)
		}
		INFO("Added missing tag column")
	}
	if !file {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN file VARCHAR NOT NULL DEFAULT "NO_FILE";`)
		if err != nil {
			panic(err)
		}
		INFO("Added missing file column")
	}
	if !line {
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN line INTEGER NOT NULL DEFAULT 1;`)
		if err != nil {
			panic(err)
		}
		INFO("Added missing line column")
	}
	Default_null_priorities(defaultPriority)
}

// TODO: get rid of this function?
// This function will be more and more useless as time goes on, but since when
// The priority column was first added it was able to have null values, this
// function might be required sometimes, so it has to be in the fix subcommand
func Default_null_priorities(defaultPriority int) {
	todoDB := open_todo_db()
	defer todoDB.Close()

	INFO("Setting possible null priority values to default value")
	sqlStatement := `UPDATE todo SET priority = $1 WHERE priority IS NULL;`
	_, err := todoDB.Exec(sqlStatement, defaultPriority)
	if err != nil {
		ERROR("Failed to edit todo content")
		panic(err)
	}
}

func Usage_fix() {
	fmt.Print(`Usage: cldl fix [-h | --help]
    Use 'cldl fix --help' to see more
`)
}

const HelpFix = `Help for cldl fix:
    Available arguments:
        --help, -h   | Show this message

    This command will be useful after making breaking
    changes to the program and its databases. It checks
    whether a local todo list has all the required columns,
    and adds them if they don't exist.`
