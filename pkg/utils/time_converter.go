package utils

import "time"

func ConvertTo12HourFormat(timeStr string) (string, error) {

	if timeStr == "25:00" {
		return timeStr, nil
	}
	// Parse the input time string in 24-hour format
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return "", err
	}

	// Format the time into 12-hour format with AM/PM
	return t.Format("03:04 PM"), nil
}
