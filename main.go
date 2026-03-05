package main

import (
	"fmt"
	"local/flagger/flagger"
	"os"
	"strings"
	"todo/cmd"

	"github.com/BurntSushi/toml"
)

func main() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		errout("Failed to get config")
		return
	}

	conf := cmd.DefaultConfig()
	configFile := configDir + "/todo/config.toml"
	_, err = toml.DecodeFile(configFile, &conf)
	if err != nil {
		panic(err)
	}

	var flags flagger.Flagset
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
		flags.Flags = []string{
			"h",
			"help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageInit()
			os.Exit(1)
		}

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "h", "help":
				cmd.HelpInit()
				return
			}
		}

		cmd.InitTodo()

	case "add":
		flags.Flags = []string{
			"h",
			"help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageAdd()
			os.Exit(1)
		}

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "h", "help":
				cmd.HelpAdd()
				return
			}
		}

		title := strings.Join(parsedArgs.NormalStr, " ")
		cmd.AddTodo(title, conf.Auto_init)

	case "list", "ls":
		flags.Flags = []string{
			"a", "all",
			"p", "pager",
			"h", "help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageList()
			os.Exit(1)
		}

		listAll := false
		pagerList := true
		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "a", "all":
				listAll = true
			case "p", "pager":
				pagerList = false
			case "h", "help":
				cmd.HelpList()
				return
			}
		}

		cmd.ListTodo(listAll, pagerList)

	case "rm", "remove", "done":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}

		flags.Flags = []string{
			"a", "all",
			"h", "help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageRm()
			os.Exit(1)
		}

		rmAll := false
		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "a", "all":
				rmAll = true
			case "h", "help":
				cmd.HelpRm()
				return
			}
		}

		title := strings.Join(parsedArgs.NormalStr, " ")
		cmd.RmTodo(title, rmAll, conf.Ask_full_rm)

	case "edit":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}

		flags.Flags = []string{
			"k", "keep",
			"h", "help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageRm()
			os.Exit(1)
		}

		keepContent := false
		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "k", "keep":
				keepContent = true
			case "h", "help":
				cmd.HelpEdit()
				return
			default:
				errout("Bad Arguments")
				cmd.UsageEdit()
				os.Exit(1)
			}
		}

		title := strings.Join(parsedArgs.NormalStr, " ")
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
      --help          Show help message
      -h              Same as '--help'
      init            Create new todo in current dir
      list            List all todo list entries
      ls              Same as 'list'
      add             Add new entry into todo list
      rm <title>      Remove todo list entry or entire list, see 'todo rm --help'
      remove <title>  Same as 'rm'
      done <title>    Same as 'rm'

  Todo application that creates local per-directory todo-lists with sqlite
  List entry titles are case-insensitive when editing or removing them,
  so be careful naming them. Adding multiple entries with the same name
  might result to undefined behavior (maybe fixed later), and trying to
  remove one of them most likely removes all.

  If a panic error occurs, most likely something went wrong when interacting
  with the sqlite databases (although it is not the only way panics can occur)

  Usable config options:
      auto_init   = bool  | automatically initialize a new local todo if it
                          | doesn't exist, or ask [y/n] to initialize
                          | [Default: false]

      ask_full_rm = bool  | Ask to fully remove the local database when the
                          | last entry gets deleted [Default: false]

`)
}
