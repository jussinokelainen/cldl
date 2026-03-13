package main

import (
	"fmt"
	"local/flagger/flagger"
	"os"
	"strconv"
	"strings"
	"time"
	"todo/cmd"

	"github.com/BurntSushi/toml"
)

func main() {
	configDir, err := os.UserHomeDir()
	if err != nil {
		errout("Failed to get config directory")
		return
	}

	// NOTE: Any calls to functions that might interact with the MasterDB
	// have to be called after this
	cmd.CreateMasterDB()
	defer cmd.MasterDB.Close()

	conf := cmd.DefaultConfig()
	configFile := configDir + "/.config/todo/config.toml"
	_, err = toml.DecodeFile(configFile, &conf)
	if err != nil {
		errout("Failed to get configs, using defaults")
		conf = cmd.DefaultConfig()
	}

	args := os.Args[1:]
	// If no args given, print usage and exit since given nothing to do
	if len(args) < 1 {
		mainUsage()
		return
	}

	var flags flagger.Flagset

	switch args[0] {
	case "fix":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}

		flags.Flags = []string{
			"h",
			"help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
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

		cmd.FixTodoTable(conf.Default_priority)
	case "set-priority", "set-p":
		if !cmd.TodoExists() {
			errout("No todo exists in current directory!")
			return
		}

		flags.Flags = []string{
			"h",
			"help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
		parsedArgs, err := flagger.ParseFlags(args[1:], flags)
		if err != nil {
			errout("Bad Arguments")
			cmd.UsagePriority()
			os.Exit(1)
		}

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "h", "help":
				cmd.HelpPriority()
				return
			}
		}
		title := strings.Join(parsedArgs.NormalStr[1:], " ")
		newPrio, err := strconv.Atoi(parsedArgs.NormalStr[0])
		if err != nil {
			errout("Invalid number for flag priority")
			os.Exit(1)
		}
		cmd.EditPriority(title, newPrio)
	case "check":
		flags.Flags = []string{
			"h",
			"help",
			"no-confirm",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
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
				conf.Ask_rm_on_check = false
			}
		}

		cmd.CheckTodos(conf.Ask_rm_on_check)

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
			"h", "help",
		}
		flags.Valued_flags = []string{
			"p", "priority",
		}
		flags.Optional_value = []string{
			"auto-init",
		}
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
				conf.Auto_init = true
			}
		}

		for _, flag := range parsedArgs.ValueFlags {
			switch flag[0] {
			case "p", "priority":
				newPrio, err := strconv.Atoi(flag[1])
				if err != nil {
					errout("Priority value must be integer")
					return
				}
				conf.Default_priority = newPrio
				conf.Ask_priority = false
			case "auto-init":
				switch flag[1] {
				case "true":
					conf.Auto_init = true
				case "false":
					conf.Auto_init = false
				default:
					errout("Bad value for auto-init: " + flag[1])
					return
				}
			}
		}

		title := strings.Join(parsedArgs.NormalStr, " ")
		cmd.AddTodo(title, conf.Auto_init, conf.Default_priority, conf.Ask_priority)

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

		timeZoneFormatted, err := time.LoadLocation(strings.TrimSpace(conf.Timezone))
		if err != nil {
			errout("Failed to parse timezone")
			os.Exit(1)
		}
		cmd.ListTodo(listAll, pagerList, timeZoneFormatted, conf.Urgent_priority, conf.In_progress_priority, conf.Colors)

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

		for _, flag := range parsedArgs.Flags {
			switch flag {
			case "k", "keep":
				conf.Keep_on_edit = !conf.Keep_on_edit
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
		cmd.EditTodo(title, conf.Keep_on_edit)
	case "relocate":
		flags.Flags = []string{
			"h",
			"help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
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
		cmd.RelocateTodo(conf.Ask_rm_on_check)

	default:
		flags.Flags = []string{
			"h", "help",
		}
		flags.Valued_flags = nil
		flags.Optional_value = nil
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
Usage: todo <COMMAND> [<args>]
	Use todo --help to see available commands and arguments
`)
}
func mainHelp() {
	fmt.Print(`
Help for todo:
  Available commands:
      --help, -h           | Show this message
      set-priority, set-p  | Set the priority of a todo entry, usage
                           | [todo <title> -p <int>] or [todo -p <int> <title>]
      init                 | Create new todo in current dir
      check                | Check all locations saved by the program whether
                           | the list actually exists. Also checks that a local
                           | todo has the right columns
      relocate             | Add todo missing from location list
      list, ls             | List all todo list entries
      add                  | Add new entry into todo list
      rm, remove, done     | Remove todo list entry or entire list
      edit                 | Edit an existing todo entry
      fix                  | Fixes the todo table, useful after breaking changes

  For more info about subcommands, use 'todo <subcommand> --help'

  Todo application that creates local per-directory todo-lists with sqlite
  List entry titles are case-insensitive when editing or removing them,
  so be careful naming them. Adding multiple entries with the same name
  might result to undefined behavior (maybe fixed later), and trying to
  remove one of them most likely removes all.

  If a panic error occurs, most likely something went wrong when interacting
  with the sqlite databases (although it is not the only way panics can occur)

  Configuration expect a file 'UserConfigDir/todo/config.toml'
  For example on a real OS the filepath is '~/.config/todo/config.toml'.

  Usable config options:
      auto_init = bool        | automatically initialize a new local todo if it
                              | doesn't exist, or ask [y/n] to initialize
                              | [Default: false]

      ask_full_rm = bool      | Ask to fully remove the local database when the
                              | last entry gets deleted. If set to false,
                              | the database will not be removed
                              | [Default: false]

      ask_rm_on_check = bool  | Ask before removing an erroneus todo location
                              | when using 'todo check'
                              | [Default: true]

      Ask_priority            | Asks for the priority that will be set to when
                              | adding new list entry.
                              | [Default: false]

      keep_on_edit = bool     | Either keep the content and add new content
                              | with a 'New edit' tag, or override entry content
                              | [Default: false]

      timezone = string       | Set the timezone used when displaying dates,
                              | By default uses current local timezone

      default_priority        | Priority that gets set to new list entries
                              | [Default: 0]

      urgent_priority         | Priority number after which entries are
                              | considered to be urgent (title printed in orange)
                              | [Default: 10]

      in_progress_priority    | Priority number after which entries are
                              | considered to be WIP (title printed in green)
                              | This value should be set to something bigger
                              | than urgent priority
                              | [Default: 100]

      [colors]                | Set custom colors for different elements.
        default               | Colors must be a string of numbers (e.g. "123")
        urgent                | between 0-255, using the corresponding value 
        wip                   | from the 256 color table
        content
        border
        dim

`)
}
