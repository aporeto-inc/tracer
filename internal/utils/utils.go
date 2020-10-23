package utils

import (
	"bytes"
	"fmt"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Tabulate print a table from data
func Tabulate(headers []string, rows [][]string) string {

	out := &bytes.Buffer{}

	colors := make([]tablewriter.Colors, len(headers))
	for i := 0; i < len(headers); i++ {
		colors[i] = tablewriter.Color(tablewriter.FgCyanColor, tablewriter.Bold)
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader(headers)
	table.AppendBulk(rows)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetHeaderLine(true)
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetHeaderColor(colors...)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	table.Render()

	return "\n" + out.String()
}

// ParseTime is a function to parse time from to and since
func ParseTime(from string, to string, since time.Duration) (fromTime time.Time, toTime time.Time, sinceDuration time.Duration, err error) {

	// Parse time if set
	if from != "" {
		fromTime, err = time.Parse(time.RFC3339, from)
		if err != nil {
			return time.Time{}, time.Time{}, 0 * time.Second, fmt.Errorf("unable to parse from duration: %s is not a valid time", from)
		}
	}

	if to != "" {
		toTime, err = time.Parse(time.RFC3339, to)
		if err != nil {
			return time.Time{}, time.Time{}, 0 * time.Second, fmt.Errorf("unable to parse to duration: %s is not a valid time", to)
		}
	}

	// If to time is not set make it now
	if toTime.IsZero() {
		toTime = time.Now().Round(time.Second)
	}

	// If a from is set compute since
	if !fromTime.IsZero() {
		return fromTime, toTime, toTime.Sub(fromTime), nil
	}

	//Ohterwise compute the from from the duration
	return toTime.Add(-since), toTime, since, nil

}
