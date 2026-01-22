package main

import (
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func helpInit() {
	fmt.Println("Help for todo init:")
	fmt.Println(" Usage: todo init [<args>]")
	fmt.Println(" Available arguments:")
	fmt.Println("    --help: show help for todo init")
}

func usageInit() {
	fmt.Println(" Usage: todo init [<args>]")
	fmt.Println("  Use todo init --help to see arguments")
}

func createNewTodo() {
	todoPath := getDbPath()
	if _, err := os.Stat(todoPath); !os.IsNotExist(err) {
		info("Todo already exists in current directory!")
		return
	}

	todoDB := openTodoDB()
	defer todoDB.Close()

	_, err := todoDB.Exec(`CREATE TABLE todo(
		title VARCHAR UNIQUE,
		content VARCHAR,
		time INTEGER
		);`)
	if err != nil {
		error("Creating new todo failed!")
		panic(err)
	}

	// Add new todo location into master location db
	sqlStatement := `INSERT INTO locations(location) VALUES($1);`
	_, err = masterDB.Exec(sqlStatement, todoPath)
	if err != nil {
		error("Adding to master DB failed!")
		panic(err)
	}

	info("New todo created!")
}

func initTodo(args []string) {
	// Default usage is without any additional arguments
	if len(args) < 1 {
		createNewTodo()
		return
	}

	switch args[0] {
	case "--help":
		helpInit()
		return
	default:
		info("Invalid arguments")
		usageInit()
	}
}

func helpRm() {
	fmt.Println("Help for todo rm:")
	fmt.Println(" Usage: todo rm [<args>]")
	fmt.Println(" Available arguments:")
	fmt.Println("    --help | show help for todo init")
	fmt.Println("    --all  | fully remove todo from current dir")
}

func usageRm() {
	fmt.Println(" Usage: todo rm [<args>]")
	fmt.Println("  Use todo rm --help to see arguments")
}

func removeAllData() {
	todoPath := getDbPath()
	if _, err := os.Stat(todoPath); os.IsNotExist(err) {
		info("No todo exists in current directory!")
		return
	}

	sqlStatement := `DELETE FROM locations WHERE location = ?;`
	_, err := masterDB.Exec(sqlStatement, todoPath)
	if err != nil {
		error("Error removing from master db")
		panic(err)
	}
	err = os.Remove(todoPath)
	if err != nil {
		error("Failed removing local db")
		panic(err)
	}

	info("Succesfully removed database!")
}

func rmTodo(args []string) {
	if len(args) < 1 {
		usageRm()
		return
	}

	switch args[0] {
	case "--help":
		helpRm()
		return
	case "--all":
		removeAllData()
		return
	default:
		title := args[0]
		todoDB := openTodoDB()
		defer todoDB.Close()

		sqlStatement := `DELETE FROM todo WHERE title= ?;`
		res, err := todoDB.Exec(sqlStatement, title)
		if err != nil {
			error("Error removing entry")
			panic(err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			error("Error getting affected rows")
			panic(err)
		}

		if rowsAffected == 0 {
			info("No entry found with title " + title)
		} else {
			info("Succesfully removed entry " + title)
		}
	}
}
