package main

import (
	"cldl/cmd"
	"cldl/flagger"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		cmd.ERROR("Failed to get config directory")
		return
	}

	cmd.Create_master_db()
	defer cmd.MasterDB.Close()

	conf := cmd.Default_config()
	configFile := homeDir + "/.config/cldl/config.toml"
	_, err = toml.DecodeFile(configFile, &conf)
	if err != nil {
		cmd.INFO("No config file found.")
		conf = cmd.Default_config()
	}

	handle_parsing(conf)
}

/*
Handles the main parsing of flags, arguments and subcommands, and then calls
required subcommand functions.
*/
func handle_parsing(conf cmd.Config) {

	args := os.Args[1:]
	// If no args given, print usage and exit since there is nothing to do
	if len(args) < 1 {
		main_usage()
		return
	}

	// Initialize flags, all subcommands will have --help and -h, so they
	// can be added here
	var flags flagger.Flagset
	flags.Flags = []string{"h", "help"}
	flags.Valued_flags = []string{}
	flags.Optional_value = []string{}

	switch args[0] {
	case "generate-configs":
		generate_configs()
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
			main_usage()
			os.Exit(1)
		}

		for _, flag := range mainFlags.Flags {
			switch flag {
			case "h", "help":
				main_help()
				return
			}
		}
	}
}

func generate_configs() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		cmd.ERROR("Failed to get config directory")
		return
	}

	configDir := homeDir + "/.config/cldl/"
	configFile := configDir + "config.toml"
	if cmd.File_exists(configFile) {
		cmd.INFO("Config file already exists, nothing to do.")
		return
	} else {
		// CLDL-ENTRY: title: macOS filepath, priority: 9, tag: pkg
		default_config_file := "/usr/share/cldl/default_config.toml"
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			cmd.ERROR("Error creating config directory.")
			os.Exit(1)
		}
		copy_file(default_config_file, configFile)
	}
}

func copy_file(src string, dest string) {
	in, err := os.Open(src)
	if err != nil {
		cmd.ERROR("Default config file not found.")
		os.Exit(1)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		cmd.ERROR("Error in creating config file.")
		os.Exit(1)
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		cmd.ERROR("Error copying data to config file")
		os.Exit(1)
	}

	os.Chmod(dest, os.FileMode(0644))
	err = out.Sync()
	if err != nil {
		panic(err)
	} else {
		cmd.OK("Successfully generated config file.")
	}
}

func handle_rename(args []string, flags flagger.Flagset, color_conf cmd.ColorConf) {
	if !cmd.Todo_exists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}

	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.Usage_rename()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.Print_help_msg(cmd.HelpRename)
			return
		}
	}

	title := strings.Join(parsedArgs.NormalStr, " ")
	cmd.Rename_todo(title, color_conf)
}

func handle_fix(args []string, flags flagger.Flagset, prio_conf cmd.PriorityConf) {
	if !cmd.Todo_exists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}

	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.Usage_fix()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.Print_help_msg(cmd.HelpFix)
			return
		}
	}

	cmd.Fix_todo_table(prio_conf.Default)
}

func handle_set(args []string, flags flagger.Flagset) {
	if !cmd.Todo_exists() {
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
		cmd.Usage_set()
		os.Exit(1)
	}
	if len(parsedArgs.Flags) < 1 && len(parsedArgs.ValueFlags) < 1 {
		cmd.Usage_set()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.Print_help_msg(cmd.HelpSet)
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
			cmd.Edit_priority(title, newPrio)
		case "t", "tag":
			title := strings.Join(parsedArgs.NormalStr, " ")
			cmd.Set_tag_to_entry(title, flag[1])
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
		cmd.Usage_check()
		os.Exit(1)
	}

	check_directories := false
	verbose_check := false

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.Print_help_msg(cmd.HelpCheck)
			return
		case "no-confirm":
			general_conf.Ask_rm_on_check = false
		case "directories", "d":
			check_directories = true
		case "verbose", "v":
			verbose_check = true
		}
	}

	cmd.Check_todos(
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
		cmd.Usage_init()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.Print_help_msg(cmd.HelpInit)
			return
		}
	}

	cmd.Init_todo()
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
		cmd.Usage_add()
		os.Exit(1)
	}

	no_ask_content := false
	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.Print_help_msg(cmd.HelpAdd)
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
			conf.Add.Ask_tags = false
		}
	}

	data.Priority = conf.Priority.Default
	data.Tag = tag
	data.Empty_content = no_ask_content

	title := strings.Join(parsedArgs.NormalStr, " ")
	cmd.Add_todo(title, conf.Add, data, conf.Colors)
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
		cmd.Usage_list()
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
			cmd.Print_help_msg(cmd.HelpList)
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
	cmd.List_todo(listAll, pagerList, conf, filterByTag, tag)
}

func handle_remove(args []string, flags flagger.Flagset, rm_conf cmd.RmConf) {
	if !cmd.Todo_exists() {
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
		cmd.Usage_rm()
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
			cmd.Print_help_msg(cmd.HelpRm)
			return
		}
	}

	title := strings.Join(parsedArgs.NormalStr, " ")
	cmd.Rm_todo(title, rmAll, rmTag, rm_conf)
}

func handle_edit(args []string, flags flagger.Flagset, conf cmd.Config) {
	if !cmd.Todo_exists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}

	additionalFlags := []string{"k", "keep"}
	flags.Flags = append(flags.Flags, additionalFlags...)
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.Usage_rm()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "k", "keep":
			conf.Edit.Keep_content = !conf.Edit.Keep_content
		case "h", "help":
			cmd.Print_help_msg(cmd.HelpEdit)
			return
		default:
			cmd.ERROR("Bad Arguments")
			cmd.Usage_edit()
			os.Exit(1)
		}
	}

	title := strings.Join(parsedArgs.NormalStr, " ")
	cmd.Edit_todo(title, conf.Edit, conf.Colors)
}

func handle_relocate(args []string, flags flagger.Flagset, general_conf cmd.GeneralConf) {
	parsedArgs, err := flagger.ParseFlags(args[1:], flags)
	if err != nil {
		cmd.ERROR("Bad Arguments")
		cmd.Usage_relocate()
		os.Exit(1)
	}

	for _, flag := range parsedArgs.Flags {
		switch flag {
		case "h", "help":
			cmd.Print_help_msg(cmd.HelpRelocate)
			return
		}
	}
	if !cmd.Todo_exists() {
		cmd.ERROR("No todo exists in current directory!")
		return
	}
	cmd.Relocate_todo(general_conf.Ask_rm_on_check)

}

// NOTE: Main help and usage functions
func main_usage() {
	fmt.Print(`Usage: cldl [-h | --help] <command> [<args>]
    Use cldl --help to see available commands
`)
}

func main_help() {
	const helpmsg = `Help for cldl:
  Flags:
      --help, -h           | Show this message

  Available commands:
      generate-configs     | Generate a config file with default values
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

  Configuration expects a file '~/.config/cldl/config.toml'.

  Default configs:
    [general]
      ask_rm_on_check = true
      timezone        = "Local"
      checkdirs       = []

    [add]
      auto_init    = false
      ask_priority = false
      ask_tags     = false

    [edit]
      keep_content = false

    [priority]
      default     = 0
      urgent      = 10
      in_progress = 100

    [rm]
      ask_full            = false
      always_confirm_full = true

    [colors]
      default = "#99FFFF"
      urgent  = "#FF8000"
      wip     = "#66FF66"
      content = "#FFFFFF"
      border  = "#FF99FF"
      dim     = "#404040"
      tag     = "#FFFF66"
      file    = "#99FFFF"`

	cmd.Print_help_msg(helpmsg)
}
