package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/term"
)

var maxWidth = 50
var applyPadding = true

func List_todo(listLocations bool, pager bool, config Config, filterByTag ListTag, tag string) {
	timeZone, err := time.LoadLocation(strings.TrimSpace(config.General.Timezone))
	if err != nil {
		ERROR("Failed to parse timezone")
		os.Exit(1)
	}
	urgentPrio := config.Priority.Urgent
	wipPrio := config.Priority.In_progress
	colors := config.Colors
	set_color_scheme(colors)
	printList := func(listString string, page bool) {
		if page {
			print_to_pager(listString)
		} else {
			fmt.Print(listString)
		}
	}

	if listLocations {
		// Default pagering behavior is opposite to normal listing, so use opposite values
		applyPadding = !pager
		printList(list_all_todo_locations(), !pager)
		return
	} else {
		if !Todo_exists() {
			ERROR("No todo exists in current directory!")
			return
		}

		todoSlice, err := get_todo_slice(filterByTag, tag)
		if err != nil {
			switch filterByTag {
			case ONLY:
				INFO("No entries with tag " + tag)
			case EXCEPT:
				INFO("No entries without tag " + tag)
			default:
				INFO("Todo list empty!")
			}
			return
		}

		if len(todoSlice) == 1 {
			pager = !pager
		}
		applyPadding = pager

		printList(format_list_items(todoSlice, timeZone, urgentPrio, wipPrio), pager)
		return
	}
}

// Formats the items in a given slice. Returns a string with all decorations, newlines etc
// that is ready to be printed as returned
func format_list_items(todoSlice []TodoStruct, timeZone *time.Location, urgentPrio int, wipPrio int) string {
	current_box := 1
	var listString strings.Builder
	// Top border and empty row
	pad_content_to_center(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "%s╔══%s══╗\n", borderColor, add_line(maxWidth))
	pad_content_to_center(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "%s║  %s  ║\033[0m", borderColor, add_space(maxWidth))

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
		listString.WriteString(make_title_line(row.Title, priorityColor))

		// Print content
		if row.Content != "[ EMPTY ]" {
			contentWrapped := wordwrap.WrapString(row.Content, uint(maxWidth))
			contentLines := strings.SplitSeq(contentWrapped, "\n")
			for line := range contentLines {
				if len(line)%2 != 0 && utf8.RuneCountInString(row.Title)%2 != 0 {
					line = " " + line
				}
				listString.WriteString(format_content_line(line))
			}
			pad_content_to_center(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "%s║  %s  ║\033[0m\n", borderColor, add_space(maxWidth))
		}

		// Print location, tag, timestamp and priority level
		if row.File != "NO_FILE" {
			listString.WriteString(make_location_line(row.File, row.Line, fileColor))
		}
		if row.Tag != "NONE" {
			listString.WriteString(make_tag_line(row.Tag))
		}
		listString.WriteString(make_priority_line(int(row.Priority), priorityColor))
		listString.WriteString(make_timestamp_line(row.Time, timeZone))

		// Print borders that continue into the next box if not last box
		if current_box < len(todoSlice) {
			pad_content_to_center(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "%s╠══%s══╣\033[0m\n", borderColor, add_line(maxWidth))
			pad_content_to_center(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "%s║  %s  ║\033[0m", borderColor, add_space(maxWidth))
		} else {
			pad_content_to_center(&listString, maxWidth+2)
			fmt.Fprintf(&listString, "%s╚══%s══╝\033[0m\n", borderColor, add_line(maxWidth))
		}
		current_box++
	}
	return listString.String()
}

func make_location_line(file_path string, file_line int, Color string) string {
	var titleStr strings.Builder
	loc_string := file_path + ":" + strconv.Itoa(file_line)

	padLeft, padRight := get_side_padding(loc_string)
	titleSize := utf8.RuneCountInString(loc_string)

	pad_content_to_center(&titleStr, maxWidth+2)
	fmt.Fprintf(&titleStr, "%s║%s", borderColor, add_space(padLeft+2))
	fmt.Fprintf(&titleStr, "%s%s", Color, loc_string)
	fmt.Fprintf(&titleStr, "%s%s║\n", borderColor, add_space(padRight+2))

	pad_content_to_center(&titleStr, maxWidth+2)
	fmt.Fprintf(&titleStr, "%s║%s", borderColor, add_space(padLeft))
	fmt.Fprintf(&titleStr, "%s══%s══", Color, add_line(titleSize))
	fmt.Fprintf(&titleStr, "%s%s║\n", borderColor, add_space(padRight))

	return titleStr.String()

}

func make_title_line(title string, priorityColor string) string {
	var titleStr strings.Builder

	padLeft, padRight := get_side_padding(title)
	titleSize := utf8.RuneCountInString(title)
	fmt.Fprint(&titleStr, "\n")

	pad_content_to_center(&titleStr, maxWidth+2)
	fmt.Fprintf(&titleStr, "%s║%s", borderColor, add_space(padLeft+2))
	fmt.Fprintf(&titleStr, "%s%s", priorityColor, title)
	fmt.Fprintf(&titleStr, "%s%s║\n", borderColor, add_space(padRight+2))

	pad_content_to_center(&titleStr, maxWidth+2)
	fmt.Fprintf(&titleStr, "%s║%s", borderColor, add_space(padLeft))
	fmt.Fprintf(&titleStr, "%s══%s══", priorityColor, add_line(titleSize))
	fmt.Fprintf(&titleStr, "%s%s║\n", borderColor, add_space(padRight))

	return titleStr.String()
}

func make_tag_line(tag string) string {
	var tagStr strings.Builder

	// Print tag
	tagString := "Tag: " + tag
	padLeft, padRight := get_side_padding(tagString)

	pad_content_to_center(&tagStr, maxWidth+2)
	fmt.Fprintf(&tagStr, "%s║%s", borderColor, add_space(padLeft+2))
	fmt.Fprintf(&tagStr, "%s%s", tagColor, tagString)
	fmt.Fprintf(&tagStr, "%s%s║\n", borderColor, add_space(padRight+2))

	return tagStr.String()
}
func make_timestamp_line(timestamp int64, timeZone *time.Location) string {
	var timeStr strings.Builder

	// Print creation timestamp
	timeString := "Created: " + time.Time.String(time.Unix(timestamp, 0).In(timeZone))
	padLeft, padRight := get_side_padding(timeString)

	pad_content_to_center(&timeStr, maxWidth+2)
	fmt.Fprintf(&timeStr, "%s║%s", borderColor, add_space(padLeft+2))
	fmt.Fprintf(&timeStr, "%s%s", dimColor, timeString)
	fmt.Fprintf(&timeStr, "%s%s║\n", borderColor, add_space(padRight+2))

	return timeStr.String()
}

func make_priority_line(priority int, priorityColor string) string {
	var prioStr strings.Builder
	if priorityColor == defaultColor {
		priorityColor = dimColor
	}
	// Print Priority
	prioString := "Priority: " + fmt.Sprint(priority)
	padLeft, padRight := get_side_padding(prioString)

	pad_content_to_center(&prioStr, maxWidth+2)
	fmt.Fprintf(&prioStr, "%s║%s", borderColor, add_space(padLeft+2))
	fmt.Fprintf(&prioStr, "%s%s", priorityColor, prioString)
	fmt.Fprintf(&prioStr, "%s%s║\n", borderColor, add_space(padRight+2))

	return prioStr.String()
}

// Formats and decorates a single line of a todo entry's content returns it as a string
func format_content_line(line string) string {
	var listString strings.Builder
	padLeft, padRight := get_side_padding(line)

	pad_content_to_center(&listString, maxWidth+2)
	fmt.Fprintf(&listString, "%s║%s", borderColor, add_space(padLeft+2))
	fmt.Fprintf(&listString, "%s%s", contentColor, line)
	fmt.Fprintf(&listString, "%s%s║\n", borderColor, add_space(padRight+2))

	return listString.String()
}

// Lists, formats and decorates all locations inside the master database
// returns a decorated string
func list_all_todo_locations() string {
	var listString strings.Builder
	var locSlice []string
	homeDir, _ := os.UserHomeDir()
	longestLoc := 0

	rows, err := MasterDB.Query(`SELECT location FROM locations;`)
	if err != nil {
		ERROR("Failed getting all locations")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var location string
		err = rows.Scan(&location)
		if err != nil {
			ERROR("Failed scanning entry content")
			panic(err)
		}
		shortenedString := strings.Replace(location, homeDir, "~", 1)
		strLen := len(shortenedString)
		if strLen > longestLoc {
			longestLoc = strLen
		}

		locSlice = append(locSlice, shortenedString)
	}
	if len(locSlice) < 1 {
		INFO("No list locations saved. 'cldl check -d' might help, see 'cldl check -h' for more.")
		os.Exit(0)
	}

	fmt.Fprint(&listString, "\n")
	pad_content_to_center(&listString, longestLoc+2)
	fmt.Fprintf(&listString, "%s%sAll todo list locations\033[0m\n",
		add_space((longestLoc/2)-9), defaultColor)
	pad_content_to_center(&listString, longestLoc+2)
	fmt.Fprintf(&listString, "%s╔══%s══╗\n", borderColor, add_line(longestLoc))

	for _, loc := range locSlice {
		lenDiff := longestLoc - len(loc)
		pad_content_to_center(&listString, longestLoc+2)
		fmt.Fprintf(&listString, "%s║%s", borderColor, add_space(lenDiff/2+2))
		fmt.Fprintf(&listString, "%s%s", contentColor, loc)
		fmt.Fprintf(&listString, "%s%s║\n", borderColor, add_space((lenDiff+1)/2+2))
	}

	pad_content_to_center(&listString, longestLoc+2)
	fmt.Fprintf(&listString, "%s╚══%s══╝\033[0m\n", borderColor, add_line(longestLoc))

	return listString.String()
}

func pad_content_to_center(listString *strings.Builder, contentWidth int) {
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
func get_todo_slice(filterByTag ListTag, tag string) ([]TodoStruct, error) {
	var todoSlice []TodoStruct
	todoDB := open_todo_db()
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
				tag,
				file,
				line
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
				tag,
				file,
				line
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
				tag,
				file,
				line
			FROM
				todo
			ORDER BY
				priority DESC,
				time ASC;
			`
	}

	rows, err := todoDB.Query(sqlStatement)
	if err != nil {
		ERROR("Failed getting todo list")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var row TodoStruct
		err = rows.Scan(&row.Title, &row.Content, &row.Time, &row.Priority, &row.Tag, &row.File, &row.Line)
		if err != nil {
			ERROR("Failed scanning entry content")
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
func get_side_padding(s string) (int, int) {
	strSize := utf8.RuneCountInString(s)
	free := max(maxWidth-strSize, 0)
	padLeft := free / 2
	padRight := free - padLeft

	return padLeft, padRight
}

// Formatting helper functions that return a string
// with a given amount of spaces or lines
func add_line(length int) string {
	if length > 0 {
		return strings.Repeat("═", length)
	} else {
		return ""
	}
}

func add_space(length int) string {
	if length > 0 {
		return strings.Repeat(" ", length)
	} else {
		return ""
	}
}

// NOTE: List command help and usage functions
func Usage_list() {
	fmt.Print(`Usage: cldl list [-h | --help] [-a | --all] [-p | --pager]
    Use 'cldl list --help' to see more
`)
}

const HelpList = `Help for cldl list:
    Available arguments:
        --help, -h   | Show this message
        --all, -a    | List all locations with todo's
        --pager, -p  | Toggle pagering behavior, by default normal lists
                     | get pagered, location lists don't
        --tag, -t    | Only list entries with a certain tag
        --except, -e | Only list entries without a certain tag

    Show content in a local todo list, or alternatively with '--all'
    show all locations with todo lists`
