package dateparser

import (
	"fmt"
	"time"
)

// ParseDate parses a date string in "YYYY-MM-DD" format and returns a time.Time object
func ParseDate(input string) (time.Time, error) {
	const layout = "2006-01-02"

	// Parse the input string
	convertedDate, err := time.Parse(layout, input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %v", err)
	}

	return convertedDate, nil
}
