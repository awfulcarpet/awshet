package main

import (
	"errors"
	"regexp"
	"time"
)

func parseCheckMessage(message string) (string, time.Time, error) {
	typeRegex := regexp.MustCompile("(in|out)")
	// TODO: implement date string
	nowRegex := regexp.MustCompile("now")

	checkType := typeRegex.FindString(message)
	if checkType == "" {
		return checkType, time.Time{}, errors.New("invalid command supplied")
	}

	date := nowRegex.FindString(message)
	if date == "" {
		return checkType, time.Time{}, errors.New("invalid time supplied")
	}

	time := time.Now()

	return checkType, time, nil
}
