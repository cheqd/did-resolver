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

	t, err := parseDateString(timeString)
	if err == nil {
		return t, nil
	}

	return time.Time{}, err
}

func parseDateString(timeString string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	// Try parsing the date using different formats
	for _, format := range formats {
		parsedTime, err := time.Parse(format, timeString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date string")
}
