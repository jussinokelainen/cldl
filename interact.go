package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func helpAdd() {
	fmt.Println("Help for todo add:")
	fmt.Println(" Available arguments:")
	fmt.Println("    --help | show help for todo add")
}

func usageAdd() {
	fmt.Println(" Usage: todo add [<args>]")
}

func addTodo(args []string) {
	if len(args) < 1 {
		usageAdd()
		return
	}

	switch args[0] {
	case "--help":
		helpAdd()
		return
	default:
		todoDB := openTodoDB()
		defer todoDB.Close()

		reader := bufio.NewReader(os.Stdin)

		title := args[0]
		time := time.Now().Unix()

		fmt.Printf("Enter contents for new todo titled %s: \n", title)
		fmt.Print("-> ")
		content, err := reader.ReadString('\n')
		if err != nil {
			error("Error reading input")
			panic(err)
		}
		content = strings.TrimSpace(content)

		sqlStatement := `INSERT INTO todo(title, content, time) VALUES($1, $2, $3);`
		_, err = todoDB.Exec(sqlStatement, title, content, time)
		if err != nil {
			error("Error adding new todo")
			info("Make sure a todo is initialized")
			panic(err)
		}
		info("Successfully added new todo " + title)
	}

}

func helpList() {
	fmt.Println("Help for todo list:")
	fmt.Println(" Available arguments:")
	fmt.Println("    --help | show help for todo list")
	fmt.Println("    -a     | list all locations with todo's")
}

func usageList() {
	fmt.Println(" Usage: todo add [<args>]")
}

func listTodo(args []string) {
	if len(args) < 1 {
		todoDB := openTodoDB()
		defer todoDB.Close()

		type TodoStruct struct {
			Title   string `json:"title"`
			Content string `json:"Content"`
			Time    int64  `json:"time"`
		}

		rows, err := todoDB.Query(`SELECT * FROM todo;`)
		if err != nil {
			error("Failed getting all locations")
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var row TodoStruct
			err = rows.Scan(&row.Title, &row.Content, &row.Time)
			if err != nil {
				error("Row scanning failed")
				panic(err)
			}
			fmt.Println(row)
		}
		err = rows.Err()
		if err != nil {
			error("Error in rows")
			panic(err)
		}
		return
	}

	switch args[0] {
	case "--help":
		helpList()
		return
	case "-a":
		listAllTodoLocations()

	default:
		usageAdd()
	}
}

func listAllTodoLocations() {
	type LocStruct struct {
		Location string `json:"location"`
	}

	rows, err := masterDB.Query(`SELECT * FROM locations;`)
	if err != nil {
		error("Failed getting all locations")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var row LocStruct
		err = rows.Scan(&row.Location)
		if err != nil {
			error("Row scanning failed")
			panic(err)
		}
		fmt.Println(row)
	}
	err = rows.Err()
	if err != nil {
		error("Error in rows")
		panic(err)
	}
}
