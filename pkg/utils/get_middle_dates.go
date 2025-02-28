package utils

import (
	"errors"
	"time"
)

func GetMiddleDates(startDate, endDate string) ([]string, error) {
	// Parse the start and end dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	// Check if start date is after end date
	if start.After(end) {
		return nil, errors.New("start date cannot be after end date")
	}

	var dates []string

	dates = append(dates, startDate)
	current := start.AddDate(0, 0, 1) // Move to the next day

	// Iterate through the middle dates
	for current.Before(end) {
		dates = append(dates, current.Format("2006-01-02"))
		current = current.AddDate(0, 0, 1) // Increment by 1 day
	}

	dates = append(dates, endDate)

	return dates, nil

}
