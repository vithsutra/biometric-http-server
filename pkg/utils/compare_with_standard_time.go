package utils

import (
	"errors"
	"time"
)

func CompareWithStandardTime(standardStartTime string, standardEndTime string, studentLoginTime string, studentLogoutTime string) (bool, error) {
	parsedStdStartTime, err := time.Parse("15:04", standardStartTime)

	if err != nil {
		return false, err
	}

	parsedStdEndTime, err := time.Parse("15:04", standardEndTime)

	if err != nil {
		return false, err
	}

	if parsedStdStartTime.After(parsedStdEndTime) {
		return false, errors.New("standard start time should be lessthan standard end time")
	}

	if parsedStdStartTime.Equal(parsedStdEndTime) {
		return false, errors.New("standard start and end time cannot be equal")
	}

	parsedStudentLoginTime, err := time.Parse("15:04", studentLoginTime)

	if err != nil {
		return false, err
	}

	parsedStudentLogoutTime, err := time.Parse("15:04", studentLogoutTime)

	if err != nil {
		return false, err
	}

	if parsedStudentLoginTime.After(parsedStudentLogoutTime) {
		return false, errors.New("student login time should be lessthan student logout time")
	}

	if parsedStudentLoginTime.Equal(parsedStudentLogoutTime) {
		return false, nil
	}

	if parsedStudentLoginTime.After(parsedStdStartTime) {
		if parsedStudentLogoutTime.Before(parsedStdEndTime) {
			return true, nil
		} else if parsedStudentLogoutTime.Equal(parsedStdEndTime) {
			return true, nil
		} else {
			return false, nil
		}
	} else if parsedStudentLoginTime.Equal(parsedStdStartTime) {
		if parsedStudentLogoutTime.Before(parsedStdEndTime) {
			return true, nil
		} else if parsedStudentLogoutTime.Equal(parsedStdEndTime) {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}

}
