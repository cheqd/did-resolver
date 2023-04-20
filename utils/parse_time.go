package utils

import "time"

func ParseFromStringTimeToGoTime(timeString string) (time.Time, error) {
	// If timeString is empty return default nullable time value (0001-01-01 00:00:00 +0000 UTC)
	if timeString == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(time.RFC3339, timeString)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(time.RFC3339Nano, timeString)
	if err == nil {
		return t, nil
	}
	return time.Time{}, err
}
