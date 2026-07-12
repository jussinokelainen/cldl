package main

import (
	"cldl/cmd"
	"cldl/flagger"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func main() {
	configDir, err := os.UserHomeDir()
	if err != nil {
		cmd.ERROR("Failed to get config directory")
		return
	}

	cmd.CreateMasterDB()
	defer cmd.MasterDB.Close()

	conf := cmd.DefaultConfig()
	configFile := configDir + "/.config/cldl/config.toml"
	_, err = toml.DecodeFile(configFile, &conf)
	if err != nil {
		cmd.ERROR("Failed to get configs, using defaults")
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
	case "init":
		handle_init(args, flags)
	case "list", "ls":
		handle_list(args, flags, conf)
	case "add":
		handle_add(args, flags, conf)
	case "rm", "remove", "done":
		handle_remove(args, flags, conf.Rm)
	case "edit":
		handle_edit(args, flags, conf)
	case "set":
		handle_set(args, flags)
	case "rename":
		handle_rename(args, flags, conf.Colors)
	case "fix":
		handle_fix(args, flags, conf.Priority)
	case "check":
		handle_check(args, flags, conf.General)
	case "relocate":
		handle_relocate(args, flags, conf.General)
	default:
		mainFlags, err := flagger.ParseFlags(args, flags)
		if err != nil || (len(mainFlags.Flags) < 1 && len(mainFlags.ValueFlags) < 1) {
			cmd.ERROR("Bad arguments")
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

func handle_rename(args []string, flags flagger.Flagset, color_conf cmd.ColorConf) {
	if !cmd.TodoExists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}

	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.UsageRename()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.PrintHelpMSG(cmd.HelpRename)
			return
		}
	}

	title := strings.Join(parsedArgs.NormalStr, " ")
	cmd.RenameTodo(title, color_conf)
}

func handle_fix(args []string, flags flagger.Flagset, prio_conf cmd.PriorityConf) {
	if !cmd.TodoExists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}

	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.UsageFix()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.PrintHelpMSG(cmd.HelpFix)
			return
		}
	}

	cmd.FixTodoTable(prio_conf.Default)
}

func handle_set(args []string, flags flagger.Flagset) {
	if !cmd.TodoExists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}

	additionalValued := []string{
		"f", "file",
		"l", "line",
		"p", "priority",
		"t", "tag",
	}
	flags.Valued_flags = append(flags.Valued_flags, additionalValued...)
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
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
			cmd.PrintHelpMSG(cmd.HelpSet)
			return
		}
	}
	for _, flag := range parsedArgs.ValueFlags {
		switch flag[0] {
		case "f", "file":
			if cmd.File_exists(flag[1]) || flag[1] == "NO_FILE" {
				title := strings.Join(parsedArgs.NormalStr, " ")
				cmd.Set_filepath_to_entry(title, flag[1])
			} else {
				cmd.ERROR("Specified file not found")
				os.Exit(1)
			}
		case "l", "line":
			linenum, err := strconv.Atoi(flag[1])
			if err != nil {
				cmd.ERROR("Invalid number for line number")
				os.Exit(1)
			}
			title := strings.Join(parsedArgs.NormalStr, " ")
			cmd.Set_fileline_to_entry(title, linenum)
		case "p", "priority":
			title := strings.Join(parsedArgs.NormalStr, " ")
			newPrio, err := strconv.Atoi(flag[1])
			if err != nil {
				cmd.ERROR("Invalid number for flag priority")
				os.Exit(1)
			}
			cmd.EditPriority(title, newPrio)
		case "t", "tag":
			title := strings.Join(parsedArgs.NormalStr, " ")
			cmd.SetTagToEntry(title, flag[1])
		}
	}
}

func handle_check(args []string, flags flagger.Flagset, general_conf cmd.GeneralConf) {
	additionalFlags := []string{
		"no-confirm",
		"directories", "d",
		"verbose", "v",
	}
	flags.Flags = append(flags.Flags, additionalFlags...)

	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.UsageCheck()
		os.Exit(1)
	}

	check_directories := false
	verbose_check := false

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.PrintHelpMSG(cmd.HelpCheck)
			return
		case "no-confirm":
			general_conf.Ask_rm_on_check = false
		case "directories", "d":
			check_directories = true
		case "verbose", "v":
			verbose_check = true
		}
	}

	cmd.CheckTodos(
		general_conf.Ask_rm_on_check,
		check_directories,
		general_conf.CheckDirs,
		verbose_check,
	)
}

func handle_init(args []string, flags flagger.Flagset) {
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.UsageInit()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.PrintHelpMSG(cmd.HelpInit)
			return
		}
	}

	cmd.InitTodo()
}

func handle_add(args []string, flags flagger.Flagset, conf cmd.Config) {
	additionalFlags := []string{"e", "empty"}
	additionalValued := []string{
		"f", "file",
		"l", "line",
		"p", "priority",
		"t", "tag",
	}
	flags.Flags = append(flags.Flags, additionalFlags...)
	flags.Valued_flags = append(flags.Valued_flags, additionalValued...)
	flags.Optional_value = append(flags.Optional_value, "auto-init")
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.UsageAdd()
		os.Exit(1)
	}

	no_ask_content := false
	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.PrintHelpMSG(cmd.HelpAdd)
			return
		case "e", "empty":
			no_ask_content = true
		case "auto-init":
			conf.Add.Auto_init = true
		}
	}

	var data cmd.AddInfo
	data.File_path = "NO_FILE"
	data.File_line = 1

	tag := "NONE"
	for _, flag := range parsedArgs.ValueFlags {
		switch flag[0] {
		case "f", "file":
			if cmd.File_exists(flag[1]) {
				data.File_path = flag[1]
			} else {
				cmd.ERROR("Specified file not found")
				os.Exit(1)
			}
		case "l", "line":
			linenum, err := strconv.Atoi(flag[1])
			if err != nil {
				cmd.ERROR("Invalid number for line number")
				os.Exit(1)
			}
			data.File_line = linenum
		case "p", "priority":
			newPrio, err := strconv.Atoi(flag[1])
			if err != nil {
				cmd.ERROR("Priority value must be integer")
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
				cmd.ERROR("Bad value for auto-init: " + flag[1])
				return
			}
		case "t", "tag":
			tag = flag[1]
		}
	}

	data.Priority = conf.Priority.Default
	data.Tag = tag
	data.Empty_content = no_ask_content

	title := strings.Join(parsedArgs.NormalStr, " ")
	cmd.AddTodo(title, conf.Add, data)
}

func handle_list(args []string, flags flagger.Flagset, conf cmd.Config) {
	additionalFlags := []string{
		"a", "all",
		"p", "pager",
	}
	additionalValued := []string{
		"t", "tag",
		"e", "except",
	}
	flags.Valued_flags = append(flags.Valued_flags, additionalValued...)
	flags.Flags = append(flags.Flags, additionalFlags...)
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
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
			cmd.PrintHelpMSG(cmd.HelpList)
			return
		}
	}

	filterByTag := cmd.ALL
	var tag string
	for _, flag := range parsedArgs.ValueFlags {
		switch flag[0] {
		case "t", "tag":
			filterByTag = cmd.ONLY
			tag = flag[1]
		case "e", "except":
			filterByTag = cmd.EXCEPT
			tag = flag[1]
		}
	}
	cmd.ListTodo(listAll, pagerList, conf, filterByTag, tag)
}

func handle_remove(args []string, flags flagger.Flagset, rm_conf cmd.RmConf) {
	if !cmd.TodoExists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}

	additionalFlags := []string{
		"a", "all",
		"t", "tag",
		"f", "file",
	}
	flags.Flags = append(flags.Flags, additionalFlags...)
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.UsageRm()
		os.Exit(1)
	}

	rmAll := false
	rmTag := false
	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "f", "file":
			title := strings.Join(parsedArgs.NormalStr, " ")
			cmd.Set_filepath_to_entry(title, "NO_FILE")
			cmd.Set_fileline_to_entry(title, 1)
			return
		case "a", "all":
			rmAll = true
		case "t", "tag":
			rmTag = true
		case "h", "help":
			cmd.PrintHelpMSG(cmd.HelpRm)
			return
		}
	}

	title := strings.Join(parsedArgs.NormalStr, " ")
	cmd.RmTodo(title, rmAll, rmTag, rm_conf)
}

func handle_edit(args []string, flags flagger.Flagset, conf cmd.Config) {
	if !cmd.TodoExists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}

	additionalFlags := []string{"k", "keep"}
	flags.Flags = append(flags.Flags, additionalFlags...)
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.UsageRm()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "k", "keep":
			conf.Edit.Keep_content = !conf.Edit.Keep_content
		case "h", "help":
			cmd.PrintHelpMSG(cmd.HelpEdit)
			return
		default:
			cmd.ERROR("Bad Arguments")
			cmd.UsageEdit()
			os.Exit(1)
		}
	}

	title := strings.Join(parsedArgs.NormalStr, " ")
	cmd.EditTodo(title, conf.Edit, conf.Colors)
}

func handle_relocate(args []string, flags flagger.Flagset, general_conf cmd.GeneralConf) {
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.UsageRelocate()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.PrintHelpMSG(cmd.HelpRelocate)
			return
		}
	}
	if !cmd.TodoExists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}
	cmd.RelocateTodo(general_conf.Ask_rm_on_check)

}

// NOTE: Main help and usage functions
func mainUsage() {
	fmt.Print(`Usage: cldl [-h | --help] <command> [<args>]
    Use cldl --help to see available commands
`)
}
func mainHelp() {
	const helpmsg = `Help for cldl:
  Flags:
      --help, -h           | Show this message

  Available commands:
      set                  | Set some values of todo entries, see cldl set --help
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
      rename               | Change the title of a todo entry

  For more info about commands, use 'cldl <command> --help'

  Todo application that creates local per-directory todo-lists with sqlite
  List entry titles are case-insensitive when editing or removing them,
  so be careful naming them. Adding multiple entries with the same name
  might result to undefined behavior (maybe fixed later), and trying to
  remove one of them most likely removes all.

  If a panic error occurs, most likely something went wrong when interacting
  with the sqlite databases (although it is not the only way panics can occur)

  Configuration expect a file '~/.config/cldl/config.toml'.

  Default configs:
    [general]
      Ask_rm_on_check = true
      Timezone = "Local"

    [add]
      Auto_init = false
      Ask_priority = false

    [edit]
      Keep_content = false

    [priority]
      Default = 0
      Urgent = 10
      In_progress = 100

    [rm]
      Ask_full = false
      Always_confirm_full = true

    [colors]
      Default = "#99FFFF"
      Urgent = "#FF8000"
      Wip = "#66FF66"
      Content = "#FFFFFF"
      Border = "#FF99FF"
      Dim = "#404040"
      Tag = "#FFFF66"`

	cmd.PrintHelpMSG(helpmsg)
}
