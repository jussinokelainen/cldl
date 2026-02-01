package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/crypto/ssh/terminal"
)

// Maximum width of the text rendered before wrapping it
var maxWidth = 50

type TodoStruct struct {
	Title   string `json:"title"`
	Content string `json:"Content"`
	Time    int64  `json:"time"`
}

func listTodo(args []string) {
	if len(args) < 1 {
		if !todoExists() {
			return
		}

		todoSlice, err := getTodoSlice()
		if err != nil {
			info("Todo list empty!")
			return
		}

		printToPager(formatListItems(todoSlice))
		return
	}

	switch args[0] {
	case "--help", "-h":
		helpList()
		return
	case "--all", "-a":
		printToPager(listAllTodoLocations())
		return
	default:
		usageList()
	}
}

// Formats the items in a given slice. Returns a string with all decorations, newlines etc
// that is ready to be printed as returned
func formatListItems(todoSlice []TodoStruct) string {
	current_box := 1
	var listString strings.Builder
	center := func() {
		fmt.Fprintf(&listString, "%s", centerListToScreen(maxWidth+2))
	}
	// Top border and empty row
	center()
	fmt.Fprintf(&listString, "\033[35m╔══%s══╗\n", addLine(maxWidth))
	fmt.Fprintf(&listString, "%s", centerListToScreen(maxWidth+2))
	fmt.Fprintf(&listString, "\033[35m║  %s  ║\033[0m", addSpace(maxWidth))

	for _, row := range todoSlice {
		// Print title
		titleSize := len(row.Title)
		titleLeftPad := maxWidth/2 - titleSize/2
		titleRightPad := maxWidth/2 - (titleSize+1)/2
		fmt.Fprint(&listString, "\n")
		center()
		fmt.Fprintf(&listString, "\033[35m║%s", addSpace(titleLeftPad+2))
		fmt.Fprintf(&listString, "\033[36m%s", row.Title)
		fmt.Fprintf(&listString, "\033[35m%s║\n", addSpace(titleRightPad+2))

		center()
		fmt.Fprintf(&listString, "\033[35m║%s", addSpace(titleLeftPad))
		fmt.Fprintf(&listString, "\033[36m══%s══", addLine(titleSize))
		fmt.Fprintf(&listString, "\033[35m%s║\n", addSpace(titleRightPad))

		// Print content
		if len(row.Content) < maxWidth {
			listString.WriteString(formatContentLine(row.Content))
		} else {
			contentWrapped := wordwrap.WrapString(row.Content, uint(maxWidth))
			contentLines := strings.SplitSeq(contentWrapped, "\n")
			for line := range contentLines {
				listString.WriteString(formatContentLine(line))
			}
		}
		center()
		fmt.Fprintf(&listString, "\033[35m║  %s  ║\033[0m\n", addSpace(maxWidth))

		// Print creation timestamp
		timeString := "Created: " + time.Time.String(time.Unix(row.Time, 0).UTC())
		timePadding := maxWidth/2 - len(timeString)/2
		center()
		fmt.Fprintf(&listString, "\033[35m║%s", addSpace(timePadding+2))
		fmt.Fprintf(&listString, "\033[38;5;8m%s", timeString)
		fmt.Fprintf(&listString, "%s\033[35m║\n", addSpace(timePadding+2))

		// Print borders that continue into the next box if not last box
		if current_box < len(todoSlice) {
			center()
			fmt.Fprintf(&listString, "\033[35m╠══%s══╣\033[0m\n", addLine(maxWidth))
			center()
			fmt.Fprintf(&listString, "\033[35m║  %s  ║\033[0m", addSpace(maxWidth))
		} else {
			center()
			fmt.Fprintf(&listString, "\033[35m╚══%s══╝\n", addLine(maxWidth))
		}
		current_box++
	}
	return listString.String()
}

// Formats and decorates a single line of a todo entry's content returns it as a string
func formatContentLine(line string) string {
	var listString strings.Builder
	visibleLen := len(line)
	free := max(maxWidth-visibleLen, 0)
	left := free / 2
	right := free - left

	fmt.Fprintf(&listString, "%s", centerListToScreen(maxWidth+2))
	fmt.Fprintf(&listString, "\033[35m║%s", addSpace(left+2))
	fmt.Fprintf(&listString, "\033[32m%s", line)
	fmt.Fprintf(&listString, "\033[35m%s║\n", addSpace(right+2))

	return listString.String()
}

// Lists, formats and decorates all locations inside the master database
// returns a decorated string
func listAllTodoLocations() string {
	var listString strings.Builder
	var locSlice []string

	rows, err := masterDB.Query(`SELECT * FROM locations;`)
	if err != nil {
		errout("Failed getting all locations")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var location string
		err = rows.Scan(&location)
		if err != nil {
			errout("Row scanning failed")
			panic(err)
		}
		locSlice = append(locSlice, location)
	}

	longestLoc := 0
	for _, loc := range locSlice {
		iLen := len(loc)
		if iLen > longestLoc {
			longestLoc = iLen
		}
	}

	fmt.Fprintf(&listString, "\n%s", centerListToScreen(longestLoc+2))
	fmt.Fprintf(&listString, "%s\033[36mAll todo list locations\033[0m\n",
		addSpace((longestLoc/2)-9))
	fmt.Fprintf(&listString, "%s", centerListToScreen(longestLoc+2))
	fmt.Fprintf(&listString, "\033[35m╔══%s══╗\n", addLine(longestLoc))

	for _, loc := range locSlice {
		lenDiff := longestLoc - len(loc)
		fmt.Fprintf(&listString, "%s", centerListToScreen(longestLoc+2))
		fmt.Fprintf(&listString, "\033[35m║%s", addSpace(lenDiff/2+2))
		fmt.Fprintf(&listString, "\033[36m%s", loc)
		fmt.Fprintf(&listString, "\033[35m%s║\n", addSpace((lenDiff+1)/2+2))
	}

	fmt.Fprintf(&listString, "%s", centerListToScreen(longestLoc+2))
	fmt.Fprintf(&listString, "\033[35m╚══%s══╝\033[0m\n", addLine(longestLoc))

	return listString.String()
}

// Simply prints the given content string into less or errors out
func printToPager(content string) {
	cmd := exec.Command("/usr/bin/less", "-R")
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		errout("Error in running less")
		panic(err)
	}
}

// Gets and returns the contents of a local todo database,
// returns an error if the list is empty
func getTodoSlice() ([]TodoStruct, error) {
	var todoSlice []TodoStruct
	todoDB := openTodoDB()
	defer todoDB.Close()

	rows, err := todoDB.Query(`SELECT * FROM todo;`)
	if err != nil {
		errout("Failed getting todo list")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var row TodoStruct
		err = rows.Scan(&row.Title, &row.Content, &row.Time)
		if err != nil {
			errout("Row scanning failed")
			panic(err)
		}
		todoSlice = append(todoSlice, row)
	}

	if len(todoSlice) < 1 {
		return todoSlice, fmt.Errorf("empty todo-list")
	}

	return todoSlice, nil
}

func centerListToScreen(contentWidth int) string {
	spacerSize := 1
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		terminalMid := width / 2
		contentMid := contentWidth / 2
		spacerSize = terminalMid - contentMid
	}

	return strings.Repeat(" ", spacerSize)
}

// Formatting helper functions that return a string
// with a given amount of spaces or lines
func addLine(length int) string {
	return strings.Repeat("═", length)
}
func addSpace(length int) string {
	return strings.Repeat(" ", length)
}
