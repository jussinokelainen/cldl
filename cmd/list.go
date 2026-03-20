package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/term"
)

var maxWidth = 50
var applyPadding = true

type TodoStruct struct {
	Title    string `json:"title"`
	Content  string `json:"Content"`
	Time     int64  `json:"time"`
	Priority int64  `json:"priority"`
	Tag      string `json:"tag"`
}

func ListTodo(listLocations bool, pager bool, config Config, filterByTag bool, tag string) {
	timeZone, err := time.LoadLocation(strings.TrimSpace(config.General.Timezone))
	if err != nil {
		errout("Failed to parse timezone")
		os.Exit(1)
	}
	urgentPrio := config.Priority.Urgent
	wipPrio := config.Priority.In_progress
	colors := config.Colors
	setColorScheme(colors)
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
			errout("No todo exists in current directory!")
			return
		}

		todoSlice, err := getTodoSlice(filterByTag, tag)
		if err != nil {
			if filterByTag {
				info("No entries with tag " + tag)
			} else {
				info("Todo list empty!")
			}
			return
		}

		if len(todoSlice) == 1 {
			pager = !pager
		}
		applyPadding = pager

		printList(formatListItems(todoSlice, timeZone, urgentPrio, wipPrio), pager)
		return
	}
}

// Formats the items in a given slice. Returns a string with all decorations, newlines etc
// that is ready to be printed as returned
func formatListItems(todoSlice []TodoStruct, timeZone *time.Location, urgentPrio int, wipPrio int) string {
	current_box := 1
	var listString strings.Builder
	// Top border and empty row
	padContentToCenter(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "%s╔══%s══╗\n", borderColor, addLine(maxWidth))
	padContentToCenter(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "%s║  %s  ║\033[0m", borderColor, addSpace(maxWidth))

	for _, row := range todoSlice {
		var priorityColor string
		if row.Priority >= int64(wipPrio) {
			priorityColor = wipColor
		} else if row.Priority >= int64(urgentPrio) {
			priorityColor = urgentColor
		} else {
			priorityColor = defaultColor
		}
		// Print title
		listString.WriteString(makeTitleLine(row.Title, priorityColor))

		// Print content
		if row.Content != "[ EMPTY ]" {
			contentWrapped := wordwrap.WrapString(row.Content, uint(maxWidth))
			contentLines := strings.SplitSeq(contentWrapped, "\n")
			for line := range contentLines {
				if len(line)%2 != 0 && utf8.RuneCountInString(row.Title)%2 != 0 {
					line = " " + line
				}
				listString.WriteString(formatContentLine(line))
			}
			padContentToCenter(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "%s║  %s  ║\033[0m\n", borderColor, addSpace(maxWidth))
		}

		// Print tag, timestamp and priority level
		if row.Tag != "NONE" {
			listString.WriteString(makeTagLine(row.Tag))
		}
		listString.WriteString(makePriorityLine(int(row.Priority), priorityColor))
		listString.WriteString(makeTimeStampLine(row.Time, timeZone))

		// Print borders that continue into the next box if not last box
		if current_box < len(todoSlice) {
			padContentToCenter(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "%s╠══%s══╣\033[0m\n", borderColor, addLine(maxWidth))
			padContentToCenter(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "%s║  %s  ║\033[0m", borderColor, addSpace(maxWidth))
		} else {
			padContentToCenter(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "%s╚══%s══╝\033[0m\n", borderColor, addLine(maxWidth))
		}
		current_box++
	}
	return listString.String()
}

func makeTitleLine(title string, priorityColor string) string {
	var titleStr strings.Builder

	titleSize := utf8.RuneCountInString(title)
	titleLeftPad := maxWidth/2 - titleSize/2
	titleRightPad := maxWidth/2 - (titleSize+1)/2
	fmt.Fprint(&titleStr, "\n")

	padContentToCenter(&titleStr, maxWidth+2)
	fmt.Fprintf(&titleStr, "%s║%s", borderColor, addSpace(titleLeftPad+2))
	fmt.Fprintf(&titleStr, "%s%s", priorityColor, title)
	fmt.Fprintf(&titleStr, "%s%s║\n", borderColor, addSpace(titleRightPad+2))

	padContentToCenter(&titleStr, maxWidth+2)
	fmt.Fprintf(&titleStr, "%s║%s", borderColor, addSpace(titleLeftPad))
	fmt.Fprintf(&titleStr, "%s══%s══", priorityColor, addLine(titleSize))
	fmt.Fprintf(&titleStr, "%s%s║\n", borderColor, addSpace(titleRightPad))

	return titleStr.String()
}

func makeTagLine(tag string) string {
	var tagStr strings.Builder

	// Print tag
	tagString := "Tag: " + tag
	tagStrlen := utf8.RuneCountInString(tagString)
	free := max(maxWidth-tagStrlen, 0)
	tagLeft := free / 2
	tagRight := free - tagLeft
	padContentToCenter(&tagStr, maxWidth+2)
	fmt.Fprintf(&tagStr, "%s║%s", borderColor, addSpace(tagLeft+2))
	fmt.Fprintf(&tagStr, "%s%s", tagColor, tagString)
	fmt.Fprintf(&tagStr, "%s%s║\n", borderColor, addSpace(tagRight+2))

	return tagStr.String()
}
func makeTimeStampLine(timestamp int64, timeZone *time.Location) string {
	var timeStr strings.Builder

	// Print creation timestamp
	timeString := "Created: " + time.Time.String(time.Unix(timestamp, 0).In(timeZone))
	timePadding := maxWidth/2 - len(timeString)/2
	padContentToCenter(&timeStr, maxWidth+2)
	fmt.Fprintf(&timeStr, "%s║%s", borderColor, addSpace(timePadding+2))
	fmt.Fprintf(&timeStr, "%s%s", dimColor, timeString)
	fmt.Fprintf(&timeStr, "%s%s║\n", borderColor, addSpace(timePadding+2))

	return timeStr.String()
}

func makePriorityLine(priority int, priorityColor string) string {
	var prioStr strings.Builder
	if priorityColor == defaultColor {
		priorityColor = dimColor
	}
	// Print Priority
	prioString := "Priority: " + fmt.Sprint(priority)
	visiblePrioLen := utf8.RuneCountInString(prioString)
	free := max(maxWidth-visiblePrioLen, 0)
	prioLeft := free / 2
	prioRight := free - prioLeft
	padContentToCenter(&prioStr, maxWidth+2)
	fmt.Fprintf(&prioStr, "%s║%s", borderColor, addSpace(prioLeft+2))
	fmt.Fprintf(&prioStr, "%s%s", priorityColor, prioString)
	fmt.Fprintf(&prioStr, "%s%s║\n", borderColor, addSpace(prioRight+2))

	return prioStr.String()
}

// Formats and decorates a single line of a todo entry's content returns it as a string
func formatContentLine(line string) string {
	var listString strings.Builder
	visibleLen := utf8.RuneCountInString(line)
	free := max(maxWidth-visibleLen, 0)
	left := free / 2
	right := free - left

	padContentToCenter(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "%s║%s", borderColor, addSpace(left+2))
	fmt.Fprintf(&listString, "%s%s", contentColor, line)
	fmt.Fprintf(&listString, "%s%s║\n", borderColor, addSpace(right+2))

	return listString.String()
}

// Lists, formats and decorates all locations inside the master database
// returns a decorated string
func listAllTodoLocations() string {
	var listString strings.Builder
	var locSlice []string
	homeDir, _ := os.UserHomeDir()
	longestLoc := 0

	rows, err := MasterDB.Query(`SELECT location FROM locations;`)
	if err != nil {
		errout("Failed getting all locations")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var location string
		err = rows.Scan(&location)
		if err != nil {
			errout("Failed scanning entry content")
			panic(err)
		}
		shortenedString := strings.Replace(location, homeDir, "~", 1)
		strLen := len(shortenedString)
		if strLen > longestLoc {
			longestLoc = strLen
		}

		locSlice = append(locSlice, shortenedString)
	}

	fmt.Fprint(&listString, "\n")
	padContentToCenter(&listString, longestLoc+2)
	fmt.Fprintf(&listString, "%s%sAll todo list locations\033[0m\n",
		addSpace((longestLoc/2)-9), defaultColor)
	padContentToCenter(&listString, longestLoc+2)
	fmt.Fprintf(&listString, "%s╔══%s══╗\n", borderColor, addLine(longestLoc))

	for _, loc := range locSlice {
		lenDiff := longestLoc - len(loc)
		padContentToCenter(&listString, longestLoc+2)
		fmt.Fprintf(&listString, "%s║%s", borderColor, addSpace(lenDiff/2+2))
		fmt.Fprintf(&listString, "%s%s", contentColor, loc)
		fmt.Fprintf(&listString, "%s%s║\n", borderColor, addSpace((lenDiff+1)/2+2))
	}

	padContentToCenter(&listString, longestLoc+2)
	fmt.Fprintf(&listString, "%s╚══%s══╝\033[0m\n", borderColor, addLine(longestLoc))

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
		errout("Error in running less, is it installed?")
		os.Exit(1)
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

		if spacerSize < 0 {
			spacerSize = 0
		}
		spacerString := strings.Repeat(" ", spacerSize)
		fmt.Fprintf(listString, "%s", spacerString)
	}
}

// Gets and returns the contents of a local todo database,
// returns an error if the list is empty
func getTodoSlice(filterByTag bool, tag string) ([]TodoStruct, error) {
	var todoSlice []TodoStruct
	todoDB := openTodoDB()
	defer todoDB.Close()

	var sqlStatement string
	if filterByTag {
		sqlStatement = fmt.Sprintf(`SELECT title, content, time, priority, tag FROM todo WHERE tag = '%s' ORDER BY priority DESC, time ASC;`, tag)
	} else {
		sqlStatement = `SELECT title, content, time, priority, tag FROM todo ORDER BY priority DESC, time ASC;`
	}

	rows, err := todoDB.Query(sqlStatement)
	if err != nil {
		errout("Failed getting todo list")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var row TodoStruct
		err = rows.Scan(&row.Title, &row.Content, &row.Time, &row.Priority, &row.Tag)
		if err != nil {
			errout("Failed scanning entry content")
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
	if length > 0 {
		return strings.Repeat("═", length)
	} else {
		return ""
	}
}

func addSpace(length int) string {
	if length > 0 {
		return strings.Repeat(" ", length)
	} else {
		return ""
	}
}

// NOTE: List command help and usage functions
func UsageList() {
	fmt.Print(`
Usage: todo list [-h | --help] [-a | --all] [-p | --pager]
    Use 'todo list --help' to see more
`)
}
func HelpList() {
	fmt.Print(`
Help for todo list:
    Available arguments:
        --help, -h   | Show this message
        --all, -a    | List all locations with todo's
        --pager, -p  | Toggle pagering behavior, by default normal lists
                     | get pagered, location lists don't

    Show content in a local todo list, or alternatively with '--all'
    show all locations with todo lists
`)
}
