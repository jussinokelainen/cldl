package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/term"
)

var maxWidth = 50
var applyPadding = true

type TodoStruct struct {
	Title   string `json:"title"`
	Content string `json:"Content"`
	Time    int64  `json:"time"`
}

func ListTodo(listLocations bool, pager bool) {
	printList := func(listString string, page bool) {
		if page {
			printToPager(listString)
		} else {
			fmt.Print(listString)
		}
	}

	if listLocations {
		// Default pagering behavior is opposite to normal listing, so use opposite values
		applyPadding = !pager
		printList(listAllTodoLocations(), !pager)
		return
	} else {
		if !TodoExists() {
			return
		}
		applyPadding = pager

		todoSlice, err := getTodoSlice()
		if err != nil {
			info("Todo list empty!")
			return
		}

		printList(formatListItems(todoSlice), pager)
		return
	}
}

// Formats the items in a given slice. Returns a string with all decorations, newlines etc
// that is ready to be printed as returned
func formatListItems(todoSlice []TodoStruct) string {
	current_box := 1
	var listString strings.Builder
	// Top border and empty row
	padContentToCenter(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "\033[35m╔══%s══╗\n", addLine(maxWidth))
	padContentToCenter(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "\033[35m║  %s  ║\033[0m", addSpace(maxWidth))

	for _, row := range todoSlice {
		// Print title
		titleSize := len(row.Title)
		titleLeftPad := maxWidth/2 - titleSize/2
		titleRightPad := maxWidth/2 - (titleSize+1)/2
		fmt.Fprint(&listString, "\n")
		padContentToCenter(&listString, maxWidth+2)
		fmt.Fprintf(&listString, "\033[35m║%s", addSpace(titleLeftPad+2))
		fmt.Fprintf(&listString, "\033[36m%s", row.Title)
		fmt.Fprintf(&listString, "\033[35m%s║\n", addSpace(titleRightPad+2))

		padContentToCenter(&listString, maxWidth+2)
		fmt.Fprintf(&listString, "\033[35m║%s", addSpace(titleLeftPad))
		fmt.Fprintf(&listString, "\033[36m══%s══", addLine(titleSize))
		fmt.Fprintf(&listString, "\033[35m%s║\n", addSpace(titleRightPad))

		// Print content
		contentWrapped := wordwrap.WrapString(row.Content, uint(maxWidth))
		contentLines := strings.SplitSeq(contentWrapped, "\n")
		for line := range contentLines {
			listString.WriteString(formatContentLine(line))
		}
		padContentToCenter(&listString, maxWidth+2)
		fmt.Fprintf(&listString, "\033[35m║  %s  ║\033[0m\n", addSpace(maxWidth))

		// Print creation timestamp
		timeString := "Created: " + time.Time.String(time.Unix(row.Time, 0).UTC())
		timePadding := maxWidth/2 - len(timeString)/2
		padContentToCenter(&listString, maxWidth+2)
		fmt.Fprintf(&listString, "\033[35m║%s", addSpace(timePadding+2))
		fmt.Fprintf(&listString, "\033[38;5;8m%s", timeString)
		fmt.Fprintf(&listString, "%s\033[35m║\n", addSpace(timePadding+2))

		// Print borders that continue into the next box if not last box
		if current_box < len(todoSlice) {
			padContentToCenter(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "\033[35m╠══%s══╣\033[0m\n", addLine(maxWidth))
			padContentToCenter(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "\033[35m║  %s  ║\033[0m", addSpace(maxWidth))
		} else {
			padContentToCenter(&listString, maxWidth+2)
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

	padContentToCenter(&listString, maxWidth+2)
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

	rows, err := MasterDB.Query(`SELECT * FROM locations;`)
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

	fmt.Fprint(&listString, "\n")
	padContentToCenter(&listString, longestLoc+2)
	fmt.Fprintf(&listString, "%s\033[36mAll todo list locations\033[0m\n",
		addSpace((longestLoc/2)-9))
	padContentToCenter(&listString, longestLoc+2)
	fmt.Fprintf(&listString, "\033[35m╔══%s══╗\n", addLine(longestLoc))

	for _, loc := range locSlice {
		lenDiff := longestLoc - len(loc)
		padContentToCenter(&listString, longestLoc+2)
		fmt.Fprintf(&listString, "\033[35m║%s", addSpace(lenDiff/2+2))
		fmt.Fprintf(&listString, "\033[36m%s", loc)
		fmt.Fprintf(&listString, "\033[35m%s║\n", addSpace((lenDiff+1)/2+2))
	}

	padContentToCenter(&listString, longestLoc+2)
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

func padContentToCenter(listString *strings.Builder, contentWidth int) {
	if applyPadding {
		spacerSize := 1
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			terminalMid := width / 2
			contentMid := contentWidth / 2
			spacerSize = terminalMid - contentMid
		}

		spacerString := strings.Repeat(" ", spacerSize)
		fmt.Fprintf(listString, "%s", spacerString)
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

// Formatting helper functions that return a string
// with a given amount of spaces or lines
func addLine(length int) string {
	return strings.Repeat("═", length)
}
func addSpace(length int) string {
	return strings.Repeat(" ", length)
}
