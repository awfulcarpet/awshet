package main

import (
	"errors"
	"regexp"
	"time"
)

func parseCheckMessage(message string) (time.Time, error) {
	nowRegex := regexp.MustCompile("now")

	date := nowRegex.FindString(message)
	if date == "" {
		return time.Time{}, errors.New("invalid time supplied")
	}

	time := time.Now()

	return time, nil
}
