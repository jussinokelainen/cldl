package main

import (
	"flag"
	"fmt"
	"local/flagger/flagger"
	"os"
	"strings"
	"todo/cmd"
)

func main() {
	args := os.Args[1:]
	// If no args given, print usage and exit
	if len(args) < 1 {
		mainUsage()
		return
	}
	cmd.CreateMasterDB()
	defer cmd.MasterDB.Close()

	switch args[0] {
	case "--help", "-h":
		mainHelp()
	case "init":
		initFlags := flag.NewFlagSet("initFlags", flag.ExitOnError)
		help := initFlags.Bool("h", false, "show help for todo init")
		helpLong := initFlags.Bool("help", false, "show help for todo init")
		initFlags.Usage = cmd.UsageInit
		err := initFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			cmd.HelpInit()
			return
		}
		if len(initFlags.Args()) > 0 {
			fmt.Println("Bad Arguments")
			initFlags.Usage()
			os.Exit(1)
		}
		cmd.InitTodo()

	case "add":
		addFlags := flag.NewFlagSet("addFlags", flag.ExitOnError)
		help := addFlags.Bool("h", false, "show help for todo add")
		helpLong := addFlags.Bool("help", false, "show help for todo add")
		addFlags.Usage = cmd.UsageAdd
		err := addFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			cmd.HelpAdd()
			return
		}
		title := strings.Join(addFlags.Args(), " ")
		cmd.AddTodo(title)

	case "list":
		listFlags := flag.NewFlagSet("listFlags", flag.ExitOnError)
		help := listFlags.Bool("h", false, "show help for todo list")
		helpLong := listFlags.Bool("help", false, "show help for todo list")
		pagerLong := listFlags.Bool("pager", false, "Do not send local list to pager")
		pager := listFlags.Bool("p", false, "Do not send local list to pager")
		all := listFlags.Bool("a", false, "List all todo list locations")
		allLong := listFlags.Bool("all", false, "List all todo list locations")
		listFlags.Usage = cmd.UsageList
		err := listFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			cmd.HelpList()
			return
		}
		listAll := false
		if *all || *allLong {
			listAll = true
		}

		pagerList := true
		if *pager || *pagerLong {
			pagerList = false
		}
		// Reverse value since bool flags don't automatically toggle to false
		cmd.ListTodo(listAll, pagerList)

	case "rm", "remove", "done":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}
		rmFlags := flag.NewFlagSet("rmFlags", flag.ExitOnError)
		help := rmFlags.Bool("h", false, "show help for todo rm")
		helpLong := rmFlags.Bool("help", false, "show help for todo rm")
		all := rmFlags.Bool("a", false, "List all todo list locations")
		allLong := rmFlags.Bool("all", false, "List all todo list locations")
		rmFlags.Usage = cmd.UsageRm
		err := rmFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			cmd.HelpRm()
			return
		}
		rmAll := false
		if *all || *allLong {
			rmAll = true
		}
		title := strings.Join(rmFlags.Args(), " ")
		cmd.RmTodo(title, rmAll)

	case "edit":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}
		editFlags := flag.NewFlagSet("editFlags", flag.ExitOnError)
		help := editFlags.Bool("h", false, "show help for todo edit")
		helpLong := editFlags.Bool("help", false, "show help for todo edit")
		keep := editFlags.Bool("k", false, "List all todo list locations")
		keepLong := editFlags.Bool("keep", false, "List all todo list locations")
		editFlags.Usage = cmd.UsageEdit
		err := editFlags.Parse(args[1:])
		if err != nil {
			os.Exit(1)
		}
		if *help || *helpLong {
			cmd.HelpEdit()
			return
		}
		keepContent := false
		if *keep || *keepLong {
			keepContent = true
		}
		title := strings.Join(editFlags.Args(), " ")
		cmd.EditTodo(title, keepContent)

	default:
		errout("Bad arguments")
		mainUsage()
	}
}

// Status printing helpers
func ok(msg string)     { fmt.Println("[\033[32m OK \033[0m] ", msg) }
func info(msg string)   { fmt.Println("[\033[35m INFO \033[0m] ", msg) }
func errout(msg string) { fmt.Println("[\033[31m ERROR \033[0m] ", msg) }

// NOTE: Main help and usage functions
func mainUsage() {
	fmt.Print(`
Usage: todo <COMMAND> [<args>]
	Use todo --help to see available commands and arguments
`)
}
func mainHelp() {
	fmt.Print(`
Help for todo:
	Available commands:
		--help		   | Show help message
		-h			   | Same as '--help'
		init		   | Create new todo in current dir
		list		   | List all todo list entries
		add			   | Add new entry into todo list
		rm <title>	   | Remove todo list entry or entire list, see 'todo rm --help'
		remove <title> | Same as 'rm'
		done <title>   | Same as 'rm'

	Todo application that creates local per-directory todo-lists with sqlite
	List entry titles are case-insensitive when editing or removing them,
	so be careful naming them. Adding multiple entries with the same name
	might result to undefined behavior (maybe fixed later), and trying to
	remove one of them most likely removes all.

	If a panic error occurs, most likely something went wrong when interacting
	with the sqlite databases (although it is not the only way panics can occur)
`)
}
