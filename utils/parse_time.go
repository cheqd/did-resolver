package utils

import (
	"fmt"
	"time"
)

func ParseFromStringTimeToGoTime(timeString string) (time.Time, error) {
	// If timeString is empty return default nullable time value (0001-01-01 00:00:00 +0000 UTC)
	if timeString == "" {
		return time.Time{}, nil
	}

	t, err := parseTimeString(timeString)
	if err == nil {
		return t, nil
	}

	return time.Time{}, err
}

func parseTimeString(timeString string) (time.Time, error) {
	formats := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.DateTime,
		time.DateOnly,
	}

	// Try parsing the date using different formats
	for _, format := range formats {
		parsedTime, err := time.Parse(format, timeString)
		if err == nil {
			return parsedTime.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date string")
}
