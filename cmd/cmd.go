package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"golang.org/x/term"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var lineComments = map[string]string{
	".go":    "//",
	".c":     "//",
	".cpp":   "//",
	".h":     "//",
	".hpp":   "//",
	".rs":    "//",
	".js":    "//",
	".ts":    "//",
	".java":  "//",
	".kt":    "//",
	".swift": "//",

	".py":   "#",
	".sh":   "#",
	".bash": "#",
	".zsh":  "#",
	".yaml": "#",
	".yml":  "#",
	".toml": "#",
	".rb":   "#",

	".lua": "--",
	".sql": "--",

	".tex": "%",
	".vim": "\"",
}

type ListTag int

const (
	ALL ListTag = iota
	ONLY
	EXCEPT
)

var MasterDB *sql.DB

var (
	defaultColor string
	urgentColor  string
	wipColor     string

	contentColor string
	borderColor  string

	dimColor string
	tagColor string
)

type TodoStruct struct {
	Title    string `json:"title"`
	Content  string `json:"Content"`
	Time     int64  `json:"time"`
	Priority int64  `json:"priority"`
	Tag      string `json:"tag"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

/*
Struct to hold all config options for this application.

Auto_init        |  Supposed to be given to AddTodo function only
Ask_full_rm      |  Supposed to be given to RmTodo function only
Ask_rm_on_check  |  Supposed to be given to CheckTodo and RelocateTodo functions
Ask_priority     |  Whether priority is asked for when adding list entry
Keep_on_edit     |  Supposed to be given to EditTodo function only
Timezone         |  Has to be formatted to *time.Location before usage

Default_priority      |  Priority that gets set to new entries, and null values
Urgent_priority       |  Priority value after which a list item is considered urgent
In_progress_priority  |  Priority value after which a list item is considered in progress

Colors  |  Set colorscheme
*/
type Config struct {
	General  GeneralConf
	Add      AddConf
	Edit     EditConf
	Priority PriorityConf
	Rm       RmConf
	Colors   ColorConf
}

type GeneralConf struct {
	Ask_rm_on_check bool
	Timezone        string
	CheckDirs       []string
}

type AddConf struct {
	Auto_init    bool
	Ask_priority bool
	Ask_tags     bool
}

type EditConf struct {
	Keep_content bool
}

type PriorityConf struct {
	Default     int
	Urgent      int
	In_progress int
}

type RmConf struct {
	Ask_full            bool
	Always_confirm_full bool
}

type ColorConf struct {
	Default string
	Urgent  string
	Wip     string
	Content string
	Border  string
	Dim     string
	Tag     string
}

func DefaultConfig() Config {
	var conf Config
	var general GeneralConf
	general.Ask_rm_on_check = true
	general.Timezone = "Local"
	general.CheckDirs = []string{}
	conf.General = general

	var add AddConf
	add.Auto_init = false
	add.Ask_priority = false
	add.Ask_tags = false
	conf.Add = add

	var edit EditConf
	edit.Keep_content = false
	conf.Edit = edit

	var priority PriorityConf
	priority.Default = 0
	priority.Urgent = 10
	priority.In_progress = 100
	conf.Priority = priority

	var rm RmConf
	rm.Ask_full = false
	rm.Always_confirm_full = true
	conf.Rm = rm

	var colors ColorConf
	colors.Default = "#99FFFF"
	colors.Urgent = "#FF8000"
	colors.Wip = "#66FF66"
	colors.Content = "#FFFFFF"
	colors.Border = "#FF99FF"
	colors.Dim = "#404040"
	colors.Tag = "#FFFF66"
	conf.Colors = colors

	return conf
}

func setColorScheme(c ColorConf) {
	var err error
	defaultColor, err = hexToRgbString(c.Default)
	if err != nil {
		ERROR("Failed to parse default color")
		os.Exit(1)
	}
	urgentColor, err = hexToRgbString(c.Urgent)
	if err != nil {
		ERROR("Failed to parse urgent color")
		os.Exit(1)
	}
	wipColor, err = hexToRgbString(c.Wip)
	if err != nil {
		ERROR("Failed to parse wip color")
		os.Exit(1)
	}

	contentColor, err = hexToRgbString(c.Content)
	if err != nil {
		ERROR("Failed to parse content color")
		os.Exit(1)
	}
	borderColor, err = hexToRgbString(c.Border)
	if err != nil {
		ERROR("Failed to parse border color")
		os.Exit(1)
	}

	dimColor, err = hexToRgbString(c.Dim)
	if err != nil {
		ERROR("Failed to parse dim color")
		os.Exit(1)
	}

	tagColor, err = hexToRgbString(c.Tag)
	if err != nil {
		ERROR("Failed to parse tag color")
		os.Exit(1)
	}
}

func hexToRgbString(hex string) (string, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return "", fmt.Errorf("invalid hex color: %s", hex)
	}
	r, err := strconv.ParseInt(hex[1:3], 16, 0)
	if err != nil {
		return "", err
	}

	g, err := strconv.ParseInt(hex[3:5], 16, 0)
	if err != nil {
		return "", err
	}

	b, err := strconv.ParseInt(hex[5:7], 16, 0)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), nil
}

func PrintHelpMSG(s string) {
	if s == "" {
		return
	}
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	maxHeight := height - 4
	lineCount := strings.Count(s, "\n") + 1
	shouldPager := lineCount > maxHeight
	if shouldPager {
		printToPager(s)
		fmt.Println(s)
	} else {
		fmt.Println(s)
	}
}

// Simply prints the given content string into less or errors out
func printToPager(content string) {
	cmd := exec.Command("/usr/bin/less", "-R")
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		ERROR("Error in running less, is it installed?")
		os.Exit(1)
	}
}

// Get the path of a local todo database, returns the path as a string
func GetDbPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		ERROR("Getting db path failed!")
		panic(err)
	}
	return cwd + "/.cldl.db"
}

// Open a connection to a local database, returns a pointer to it
func openTodoDB() *sql.DB {
	db, err := sql.Open("sqlite", GetDbPath())
	if err != nil {
		ERROR("Opening todo DB failed!")
		panic(err)
	}
	return db
}

// Creates a 'master' database that holds all the locations to
// local databases if one does not yet exist
func CreateMasterDB() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		ERROR("Getting homedir failed!")
		panic(err)
	}
	masterDbDIR := homedir + "/.local/share/cldl"
	err = os.MkdirAll(masterDbDIR, 0755)
	if err != nil {
		ERROR("Creating ~/.local/share/cldl directory failed!")
		panic(err)
	}

	MasterDB, err = sql.Open("sqlite", masterDbDIR+"/.cldl.db")
	if err != nil {
		ERROR("Opening master database failed!")
		panic(err)
	}
	_, err = MasterDB.Exec(`CREATE TABLE IF NOT EXISTS locations (location VARCHAR UNIQUE);`)
	if err != nil {
		ERROR("Creating master database failed!")
		panic(err)
	}
}

// Function to add new paths for lists
// to the location database
func addToMasterDB(path string) {
	sqlStatement := `INSERT INTO locations(location) VALUES($1);`
	_, err := MasterDB.Exec(sqlStatement, path)
	if err != nil {
		ERROR("Adding to master DB failed!")
		panic(err)
	}
}

func check_if_masterdb_has_loc(location string) bool {
	sqlStatement := "SELECT location FROM locations WHERE location = ?;"
	err := MasterDB.QueryRow(sqlStatement, location).Scan(&location)
	if err != nil {
		if err != sql.ErrNoRows {
			ERROR(err)
			os.Exit(1)
		}
		return false
	}

	return true
}

func removeFromMasterDB(todoPath string) {
	sqlStatement := `DELETE FROM locations WHERE location = ?;`
	_, err := MasterDB.Exec(sqlStatement, todoPath)
	if err != nil {
		ERROR("Error removing from master db")
		panic(err)
	}
	OK("Removed \033[35m" + todoPath + "\033[0m from location list")
}

// Get the existing content of a todo entry if it exists, returns an error if it doesn't exist
func get_content_if_entry_exists(title string) (string, error) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	sqlStatement := `SELECT content from todo WHERE UPPER(title) = UPPER($1);`
	res, err := todoDB.Query(sqlStatement, title)
	if err != nil {
		ERROR("Failed checking entry in database!")
		panic(err)
	}
	defer res.Close()

	var content string
	for res.Next() {
		err = res.Scan(&content)
		if err != nil {
			ERROR("Failed scanning existing entry content")
			panic(err)
		}
	}
	if content == "" {
		return content, fmt.Errorf("No content found")
	}

	return content, nil
}

func askYesNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + " [y/n]: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		ERROR("Error reading input")
		return askYesNo(question)
	}
	answer = strings.TrimSpace(answer)
	switch answer {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Print("Invalid answer, try again.\n")
		return askYesNo(question)
	}
}

func insert_line(filename string, lineNum int, newLine string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	if lineNum < 0 {
		lineNum = 0
	}
	if lineNum > len(lines) {
		lineNum = len(lines)
	}

	lines = append(lines[:lineNum], append([]string{newLine}, lines[lineNum:]...)...)

	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}

// Checks whether a local todo exists in the current directory,
// might give an erroneus result if some other file is named
// exactly as the todo database should be,
// in which case errors that might come are a skill issue
func TodoExists() bool {
	return File_exists(GetDbPath())
}

// Check whether given file exists
func File_exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Status printing helpers
func OK(msg ...any)    { fmt.Println(append([]any{"[\033[32m OK \033[0m]"}, msg...)...) }
func INFO(msg ...any)  { fmt.Println(append([]any{"[\033[35m INFO \033[0m]"}, msg...)...) }
func ERROR(msg ...any) { fmt.Println(append([]any{"[\033[31m ERROR \033[0m]"}, msg...)...) }
