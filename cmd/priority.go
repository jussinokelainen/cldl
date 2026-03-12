package cmd

import (
	"database/sql"
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

func AddPriorityColumnIfNotExist(defaultPriority int) {
	exists := false
	info("checking for priority column")

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
		if name == "priority" {
			exists = true
			info("Priority column exists")
		}
	}

	if !exists {
		info("Added missing priority column")
		_, err := todoDB.Exec(`ALTER TABLE todo ADD COLUMN priority INTEGER`)
		if err != nil {
			panic(err)
		}
	}
	info("Setting null priority values to default value")
	DefaultNullPriorities(defaultPriority)
}

func DefaultNullPriorities(defaultPriority int) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	sqlStatement := `UPDATE todo SET priority = $1 WHERE priority IS NULL;`
	_, err := todoDB.Exec(sqlStatement, defaultPriority)
	if err != nil {
		errout("Failed to edit todo content")
		panic(err)
	}
}
