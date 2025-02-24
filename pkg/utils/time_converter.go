package utils

import "time"

func ConvertTo12HourFormat(timeStr string) (string, error) {
	// Parse the input time string in 24-hour format
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return "", err
	}

	// Format the time into 12-hour format with AM/PM
	return t.Format("03:04 PM"), nil
}
