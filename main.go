package main

import (
	"fmt"
	"local/flagger/flagger"
	"os"
	"strconv"
	"strings"
	"todo/cmd"

	"github.com/BurntSushi/toml"
)

func main() {
	configDir, err := os.UserHomeDir()
	if err != nil {
		errout("Failed to get config directory")
		return
	}

	cmd.CreateMasterDB()
	defer cmd.MasterDB.Close()

	conf := cmd.DefaultConfig()
	configFile := configDir + "/.config/todo/config.toml"
	_, err = toml.DecodeFile(configFile, &conf)
	if err != nil {
		errout("Failed to get configs, using defaults")
		conf = cmd.DefaultConfig()
	}

	handleParsing(conf)
}

/*
Handles the main parsing of flags, arguments and subcommands, and then calls
required subcommand functions.
*/
func handleParsing(conf cmd.Config) {

	args := os.Args[1:]
	// If no args given, print usage and exit since there is nothing to do
	if len(args) < 1 {
		mainUsage()
		return
	}

	// Initialize flags, all subcommands will have --help and -h, so they
	// can be added here
	var flags flagger.Flagset
	flags.Flags = []string{"h", "help"}
	flags.Valued_flags = []string{}
	flags.Optional_value = []string{}

	switch args[0] {
	case "fix":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}

		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageFix()
			os.Exit(1)
		}

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "h", "help":
				cmd.HelpFix()
				return
			}
		}

		cmd.FixTodoTable(conf.Priority.Default)
	case "set":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}

		additionalValued := []string{
			"p", "priority",
			"t", "tag",
		}
		flags.Valued_flags = append(flags.Valued_flags, additionalValued...)
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageSet()
			os.Exit(1)
		}
		if len(parsedArgs.Flags) < 1 && len(parsedArgs.ValueFlags) < 1 {
			cmd.UsageSet()
			os.Exit(1)
		}

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "h", "help":
				cmd.HelpSet()
				return
			}
		}
		for _, flag := range parsedArgs.ValueFlags {
			switch flag[0] {
			case "p", "priority":
				title := strings.Join(parsedArgs.NormalStr, " ")
				newPrio, err := strconv.Atoi(flag[1])
				if err != nil {
					errout("Invalid number for flag priority")
					os.Exit(1)
				}
				cmd.EditPriority(title, newPrio)
			case "t", "tag":
				title := strings.Join(parsedArgs.NormalStr, " ")
				cmd.SetTagToEntry(title, flag[1])
			}
		}
	case "check":
		flags.Flags = append(flags.Flags, "no-confirm")
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageCheck()
			os.Exit(1)
		}

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "h", "help":
				cmd.HelpCheck()
				return
			case "no-confirm":
				conf.General.Ask_rm_on_check = false
			}
		}

		cmd.CheckTodos(conf.General.Ask_rm_on_check)

	case "init":
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
		additionalValued := []string{
			"p", "priority",
			"t", "tag",
		}
		flags.Valued_flags = append(flags.Valued_flags, additionalValued...)
		flags.Optional_value = append(flags.Optional_value, "auto-init")
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
			case "auto-init":
				conf.Add.Auto_init = true
			}
		}

		tag := "NONE"
		for _, flag := range parsedArgs.ValueFlags {
			switch flag[0] {
			case "p", "priority":
				newPrio, err := strconv.Atoi(flag[1])
				if err != nil {
					errout("Priority value must be integer")
					return
				}
				conf.Priority.Default = newPrio
				conf.Add.Ask_priority = false
			case "auto-init":
				switch flag[1] {
				case "true":
					conf.Add.Auto_init = true
				case "false":
					conf.Add.Auto_init = false
				default:
					errout("Bad value for auto-init: " + flag[1])
					return
				}
			case "t", "tag":
				tag = flag[1]
			}
		}

		title := strings.Join(parsedArgs.NormalStr, " ")
		cmd.AddTodo(title, conf.Add, conf.Priority.Default, tag)

	case "list", "ls":
		additionalFlags := []string{
			"a", "all",
			"p", "pager",
		}
		additionalValued := []string{"t", "tag"}
		flags.Valued_flags = append(flags.Valued_flags, additionalValued...)
		flags.Flags = append(flags.Flags, additionalFlags...)
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

		filterByTag := false
		var tag string
		for _, flag := range parsedArgs.ValueFlags {
			switch flag[0] {
			case "t", "tag":
				filterByTag = true
				tag = flag[1]
			}
		}
		cmd.ListTodo(listAll, pagerList, conf, filterByTag, tag)

	case "rm", "remove", "done":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}

		additionalFlags := []string{"a", "all"}
		flags.Flags = append(flags.Flags, additionalFlags...)
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
		cmd.RmTodo(title, rmAll, conf.Rm)

	case "edit":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}

		additionalFlags := []string{"k", "keep"}
		flags.Flags = append(flags.Flags, additionalFlags...)
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageRm()
			os.Exit(1)
		}

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "k", "keep":
				conf.Edit.Keep_content = !conf.Edit.Keep_content
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
		cmd.EditTodo(title, conf.Edit, conf.Colors)
	case "relocate":
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsageRelocate()
			os.Exit(1)
		}

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "h", "help":
				cmd.HelpRelocate()
				return
			}
		}
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}
		cmd.RelocateTodo(conf.General.Ask_rm_on_check)

	default:
		mainFlags, err := flagger.ParseFlags(args, flags)
		if err != nil || (len(mainFlags.Flags) < 1 && len(mainFlags.ValueFlags) < 1) {
			errout("Bad arguments")
			mainUsage()
			os.Exit(1)
		}

		for _, flag := range mainFlags.Flags {
			switch flag {
			case "h", "help":
				mainHelp()
				return
			}
		}
	}
}

// Status printing helpers
func ok(msg string)     { fmt.Println("[\033[32m OK \033[0m] ", msg) }
func info(msg string)   { fmt.Println("[\033[35m INFO \033[0m] ", msg) }
func errout(msg string) { fmt.Println("[\033[31m ERROR \033[0m] ", msg) }

// NOTE: Main help and usage functions
func mainUsage() {
	fmt.Print(`
Usage: todo [-h | --help] <command> [<args>]
    Use todo --help to see available commands
`)
}
func mainHelp() {
	fmt.Print(`
Help for todo:
  Flags:
      --help, -h           | Show this message

  Available commands:
      set                  | Set some values of todo entries, see todo set --help
      init                 | Create new todo in current directory
      check                | Check all locations saved by the program whether
                           | the list actually exists. Also checks that a local
                           | todo has the right columns
      relocate             | Add todo missing from location list
      list, ls             | List all todo list entries
      add                  | Add new entry into todo list
      rm, remove, done     | Remove todo list entry or entire list
      edit                 | Edit an existing todo entry
      fix                  | Fixes the todo table, useful after breaking changes

  For more info about commands, use 'todo <command> --help'

  Todo application that creates local per-directory todo-lists with sqlite
  List entry titles are case-insensitive when editing or removing them,
  so be careful naming them. Adding multiple entries with the same name
  might result to undefined behavior (maybe fixed later), and trying to
  remove one of them most likely removes all.

  If a panic error occurs, most likely something went wrong when interacting
  with the sqlite databases (although it is not the only way panics can occur)

  Configuration expect a file as '~/.config/todo/config.toml'.

`)
}
