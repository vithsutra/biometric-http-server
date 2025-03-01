package utils

import (
	"errors"
	"time"
)

func GetMiddleDays(startDate, endDate string) ([]int, error) {
	const layout = "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}
	if start.After(end) {
		return nil, errors.New("start date must be before or equal to end date")
	}

	var days []int
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		days = append(days, d.Day())
	}

	return days, nil
}
