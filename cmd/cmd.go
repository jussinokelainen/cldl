package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

var MasterDB *sql.DB

var (
	defaultColor string
	urgentColor  string
	wipColor     string

	contentColor string
	borderColor  string

	dimColor string
)

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
	Auto_init       bool
	Ask_full_rm     bool
	Ask_rm_on_check bool
	Ask_priority    bool
	Keep_on_edit    bool
	Timezone        string

	Default_priority     int
	Urgent_priority      int
	In_progress_priority int

	Colors ColorConf
}

type ColorConf struct {
	Default string
	Urgent  string
	Wip     string
	Content string
	Border  string
	Dim     string
}

func DefaultConfig() Config {
	var colors ColorConf
	colors.Default = "#99FFFF"
	colors.Urgent = "#FF8000"
	colors.Wip = "#66FF66"
	colors.Content = "#FFFFFF"
	colors.Border = "#FF99FF"
	colors.Dim = "#404040"

	var conf Config
	conf.Auto_init = false
	conf.Ask_full_rm = false
	conf.Ask_rm_on_check = true
	conf.Keep_on_edit = false
	conf.Timezone = "Local"
	conf.Default_priority = 0
	conf.Urgent_priority = 10
	conf.In_progress_priority = 100
	conf.Ask_priority = false
	conf.Colors = colors

	return conf
}

func setColorScheme(c ColorConf) {
	var err error
	defaultColor, err = hexToRgbString(c.Default)
	if err != nil {
		errout("Failed to parse default color")
		os.Exit(1)
	}
	urgentColor, err = hexToRgbString(c.Urgent)
	if err != nil {
		errout("Failed to parse urgent color")
		os.Exit(1)
	}
	wipColor, err = hexToRgbString(c.Wip)
	if err != nil {
		errout("Failed to parse wip color")
		os.Exit(1)
	}

	contentColor, err = hexToRgbString(c.Content)
	if err != nil {
		errout("Failed to parse content color")
		os.Exit(1)
	}
	borderColor, err = hexToRgbString(c.Border)
	if err != nil {
		errout("Failed to parse border color")
		os.Exit(1)
	}

	dimColor, err = hexToRgbString(c.Dim)
	if err != nil {
		errout("Failed to parse dim color")
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

func addToMasterDB(path string) {
	// Add new todo location into list location database
	sqlStatement := `INSERT INTO locations(location) VALUES($1);`
	_, err := MasterDB.Exec(sqlStatement, path)
	if err != nil {
		errout("Adding to master DB failed!")
		panic(err)
	}
}

// Creates a 'master' database that holds all the locations to
// local databases if one does not yet exist
func CreateMasterDB() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		errout("Getting homedir failed!")
		panic(err)
	}
	masterDbDIR := homedir + "/.sqlite/todo"
	err = os.MkdirAll(masterDbDIR, 0755)
	if err != nil {
		errout("Creating .sqlite/todo directory failed!")
		panic(err)
	}

	MasterDB, err = sql.Open("sqlite", masterDbDIR+"/.todo.db")
	if err != nil {
		errout("Opening master database failed!")
		panic(err)
	}
	_, err = MasterDB.Exec(`CREATE TABLE IF NOT EXISTS locations (location VARCHAR UNIQUE);`)
	if err != nil {
		errout("Creating master database failed!")
		panic(err)
	}
}

// Get the path of a local todo database, returns the path as a string
func GetDbPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		errout("Getting db path failed!")
		panic(err)
	}
	return cwd + "/.todoApp.db"
}

// Open a connection to a local database, returns a pointer to it
func openTodoDB() *sql.DB {
	db, err := sql.Open("sqlite", GetDbPath())
	if err != nil {
		errout("Opening todo DB failed!")
		panic(err)
	}
	return db
}

func remove_master_entry(todoPath string) {
	sqlStatement := `DELETE FROM locations WHERE location = ?;`
	_, err := MasterDB.Exec(sqlStatement, todoPath)
	if err != nil {
		errout("Error removing from master db")
		panic(err)
	}
}

// Get the existing content of a todo entry if it exists, returns an error if it doesn't exist
func getIfEntryExists(title string) (string, error) {
	todoDB := openTodoDB()
	defer todoDB.Close()

	sqlStatement := `SELECT content from todo WHERE UPPER(title) = UPPER($1);`
	res, err := todoDB.Query(sqlStatement, title)
	if err != nil {
		errout("Failed checking entry in database!")
		panic(err)
	}
	defer res.Close()

	var content string
	for res.Next() {
		err = res.Scan(&content)
		if err != nil {
			errout("Failed scanning existing entry content")
			panic(err)
		}
	}
	if content == "" {
		return content, fmt.Errorf("No content found")
	}

	return content, nil
}

// Checks whether a local todo exists in the current directory,
// might give an erroneus result if some other file is named
// exactly as the todo database should be,
// in which case errors that might come are a skill issue
func TodoExists() bool {
	if _, err := os.Stat(GetDbPath()); os.IsNotExist(err) {
		return false
	}
	return true
}

// Status printing helpers
func ok(msg string)     { fmt.Println("[\033[32m OK \033[0m] ", msg) }
func info(msg string)   { fmt.Println("[\033[35m INFO \033[0m] ", msg) }
func errout(msg string) { fmt.Println("[\033[31m ERROR \033[0m] ", msg) }
