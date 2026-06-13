package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/term"
)

var maxWidth = 50
var applyPadding = true

func ListTodo(listLocations bool, pager bool, config Config, filterByTag ListTag, tag string) {
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
			switch filterByTag {
			case ONLY:
				info("No entries with tag " + tag)
			case EXCEPT:
				info("No entries without tag " + tag)
			default:
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

	padLeft, padRight := getSidePadding(title)
	titleSize := utf8.RuneCountInString(title)
	fmt.Fprint(&titleStr, "\n")

	padContentToCenter(&titleStr, maxWidth+2)
	fmt.Fprintf(&titleStr, "%s║%s", borderColor, addSpace(padLeft+2))
	fmt.Fprintf(&titleStr, "%s%s", priorityColor, title)
	fmt.Fprintf(&titleStr, "%s%s║\n", borderColor, addSpace(padRight+2))

	padContentToCenter(&titleStr, maxWidth+2)
	fmt.Fprintf(&titleStr, "%s║%s", borderColor, addSpace(padLeft))
	fmt.Fprintf(&titleStr, "%s══%s══", priorityColor, addLine(titleSize))
	fmt.Fprintf(&titleStr, "%s%s║\n", borderColor, addSpace(padRight))

	return titleStr.String()
}

func makeTagLine(tag string) string {
	var tagStr strings.Builder

	// Print tag
	tagString := "Tag: " + tag
	padLeft, padRight := getSidePadding(tagString)

	padContentToCenter(&tagStr, maxWidth+2)
	fmt.Fprintf(&tagStr, "%s║%s", borderColor, addSpace(padLeft+2))
	fmt.Fprintf(&tagStr, "%s%s", tagColor, tagString)
	fmt.Fprintf(&tagStr, "%s%s║\n", borderColor, addSpace(padRight+2))

	return tagStr.String()
}
func makeTimeStampLine(timestamp int64, timeZone *time.Location) string {
	var timeStr strings.Builder

	// Print creation timestamp
	timeString := "Created: " + time.Time.String(time.Unix(timestamp, 0).In(timeZone))
	padLeft, padRight := getSidePadding(timeString)

	padContentToCenter(&timeStr, maxWidth+2)
	fmt.Fprintf(&timeStr, "%s║%s", borderColor, addSpace(padLeft+2))
	fmt.Fprintf(&timeStr, "%s%s", dimColor, timeString)
	fmt.Fprintf(&timeStr, "%s%s║\n", borderColor, addSpace(padRight+2))

	return timeStr.String()
}

func makePriorityLine(priority int, priorityColor string) string {
	var prioStr strings.Builder
	if priorityColor == defaultColor {
		priorityColor = dimColor
	}
	// Print Priority
	prioString := "Priority: " + fmt.Sprint(priority)
	padLeft, padRight := getSidePadding(prioString)

	padContentToCenter(&prioStr, maxWidth+2)
	fmt.Fprintf(&prioStr, "%s║%s", borderColor, addSpace(padLeft+2))
	fmt.Fprintf(&prioStr, "%s%s", priorityColor, prioString)
	fmt.Fprintf(&prioStr, "%s%s║\n", borderColor, addSpace(padRight+2))

	return prioStr.String()
}

// Formats and decorates a single line of a todo entry's content returns it as a string
func formatContentLine(line string) string {
	var listString strings.Builder
	padLeft, padRight := getSidePadding(line)

	padContentToCenter(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "%s║%s", borderColor, addSpace(padLeft+2))
	fmt.Fprintf(&listString, "%s%s", contentColor, line)
	fmt.Fprintf(&listString, "%s%s║\n", borderColor, addSpace(padRight+2))

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
func getTodoSlice(filterByTag ListTag, tag string) ([]TodoStruct, error) {
	var todoSlice []TodoStruct
	todoDB := openTodoDB()
	defer todoDB.Close()

	var sqlStatement string
	switch filterByTag {
	case ONLY:
		sqlStatement = fmt.Sprintf(`
			SELECT
				title,
				content,
				time,
				priority,
				tag
			FROM
				todo
			WHERE
				tag = '%s'
			ORDER BY
				priority DESC,
				time ASC;
			`, tag)
	case EXCEPT:
		sqlStatement = fmt.Sprintf(`
			SELECT
				title,
				content,
				time,
				priority,
				tag
			FROM
				todo
			WHERE
				tag <> '%s'
			ORDER BY
				priority DESC,
				time ASC;
			`, tag)
	default:
		sqlStatement = `
			SELECT
				title,
				content,
				time,
				priority,
				tag
			FROM
				todo
			ORDER BY
				priority DESC,
				time ASC;
			`
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

// Helper function for getting the correct spacing in various line inside
// the printed todo list box
func getSidePadding(s string) (int, int) {
	strSize := utf8.RuneCountInString(s)
	free := max(maxWidth-strSize, 0)
	padLeft := free / 2
	padRight := free - padLeft

	return padLeft, padRight
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
	fmt.Print(`Usage: todo list [-h | --help] [-a | --all] [-p | --pager]
    Use 'todo list --help' to see more
`)
}
func HelpList() {
	const helpmsg = `Help for todo list:
    Available arguments:
        --help, -h   | Show this message
        --all, -a    | List all locations with todo's
        --pager, -p  | Toggle pagering behavior, by default normal lists
                     | get pagered, location lists don't
        --tag, -t    | Only list entries with a certain tag
        --except, -e | Only list entries without a certain tag

    Show content in a local todo list, or alternatively with '--all'
    show all locations with todo lists`

	PrintHelpMSG(helpmsg)
}
