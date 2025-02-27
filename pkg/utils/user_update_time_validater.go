package utils

import (
	"errors"
	"time"
)

func UserUpdateTimeValidater(morning_start, morning_end, afternoon_start, afternoon_end, evening_start, evening_end string) error {
	morning_start_time, err := time.Parse("15:04", morning_start)
	if err != nil {
		return err
	}
	morning_end_time, err := time.Parse("15:05", morning_end)
	if err != nil {
		return err
	}

	if !morning_end_time.After(morning_start_time) {
		return errors.New("morning_start_time should be less than morning_end_time")
	}

	afternoon_start_time, err := time.Parse("15:04", afternoon_start)

	if err != nil {
		return err
	}

	afternoon_end_time, err := time.Parse("15:04", afternoon_end)

	if err != nil {
		return err
	}

	if !afternoon_start_time.After(morning_end_time) {
		return errors.New("morning_end_time should be less than afternoon_start_time")
	}

	if !afternoon_end_time.After(afternoon_start_time) {
		return errors.New("afternoon_start_time should be less than afternoon_end_time")
	}

	evening_start_time, err := time.Parse("15:04", evening_start)

	if err != nil {
		return err
	}

	evening_end_time, err := time.Parse("15:04", evening_end)

	if err != nil {
		return err
	}

	if !evening_start_time.After(afternoon_end_time) {
		return errors.New("afternoon_end_time should be less than evening_start_time")
	}

	if !evening_end_time.After(evening_start_time) {
		return errors.New("evening_start_time should be less than evening_end_time")
	}

	return nil

}
