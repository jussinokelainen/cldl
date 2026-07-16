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
	var out_str strings.Builder
	// Top border and empty row
	pad_content_to_center(&out_str, maxWidth+2)
	fmt.Fprintf(&out_str, "%s╔══%s══╗\n", borderColor, add_line(maxWidth))
	pad_content_to_center(&out_str, maxWidth+2)
	fmt.Fprintf(&out_str, "%s║  %s  ║\033[0m", borderColor, add_space(maxWidth))

	for _, row := range todoSlice {
		var priorityColor string
		if row.Priority >= int64(wipPrio) {
			priorityColor = wipColor
		} else if row.Priority >= int64(urgentPrio) {
			priorityColor = urgentColor
		} else {
			priorityColor = defaultColor
		}
		fmt.Fprint(&out_str, "\n")
		// Print title
		titleSize := utf8.RuneCountInString(row.Title)
		fmt.Fprint(&out_str, format_line(row.Title, &priorityColor))
		fmt.Fprint(&out_str, format_line(add_line(titleSize+4), &priorityColor))

		// Print content
		if row.Content != "[ EMPTY ]" {
			contentWrapped := wordwrap.WrapString(row.Content, uint(maxWidth))
			contentLines := strings.SplitSeq(contentWrapped, "\n")
			for line := range contentLines {
				if len(line)%2 != 0 && utf8.RuneCountInString(row.Title)%2 != 0 {
					line = " " + line
				}
				fmt.Fprint(&out_str, format_line(line, &contentColor))
			}
			pad_content_to_center(&out_str, maxWidth+2)
			fmt.Fprintf(&out_str, "%s║  %s  ║\033[0m\n", borderColor, add_space(maxWidth))
		}

		// Print location, tag, timestamp and priority level
		if row.File != "NO_FILE" {
			loc_string := row.File + ":" + strconv.Itoa(row.Line)
			titleSize := utf8.RuneCountInString(loc_string)

			fmt.Fprint(&out_str, format_line(loc_string, &fileColor))
			fmt.Fprint(&out_str, format_line(add_line(titleSize+4), &fileColor))
		}
		if row.Tag != "NONE" {
			tagString := "Tag: " + row.Tag
			fmt.Fprint(&out_str, format_line(tagString, &tagColor))
		}
		// Print Priority
		if priorityColor == defaultColor {
			priorityColor = dimColor
		}
		prioString := "Priority: " + fmt.Sprint(int(row.Priority))
		fmt.Fprint(&out_str, format_line(prioString, &priorityColor))
		// Print timestamp
		timeString := "Created: " + time.Time.String(time.Unix(row.Time, 0).In(timeZone))
		fmt.Fprint(&out_str, format_line(timeString, &dimColor))

		// Print borders that continue into the next box if not last box
		if current_box < len(todoSlice) {
			pad_content_to_center(&out_str, maxWidth+2)
			fmt.Fprintf(&out_str, "%s╠══%s══╣\033[0m\n", borderColor, add_line(maxWidth))
			pad_content_to_center(&out_str, maxWidth+2)
			fmt.Fprintf(&out_str, "%s║  %s  ║\033[0m", borderColor, add_space(maxWidth))
		} else {
			pad_content_to_center(&out_str, maxWidth+2)
			fmt.Fprintf(&out_str, "%s╚══%s══╝\033[0m\n", borderColor, add_line(maxWidth))
		}
		current_box++
	}
	return out_str.String()
}

// Formats given content, adding borders and correct padding
// to center the content inside the borders
func format_line(content string, color *string) string {
	var out_str strings.Builder

	left_pad, right_pad := get_side_padding(content)
	pad_content_to_center(&out_str, maxWidth+2)
	fmt.Fprintf(&out_str, "%s║%s", borderColor, add_space(left_pad+2))
	fmt.Fprintf(&out_str, "%s%s", *color, content)
	fmt.Fprintf(&out_str, "%s%s║\n", borderColor, add_space(right_pad+2))

	return out_str.String()
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
