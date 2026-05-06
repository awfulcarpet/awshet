package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseTime(str string) (int, int, error) {
	hour_and_minute := strings.Split(str, ":")
	if len(hour_and_minute) < 2 {
		return -1, -1, errors.New("must have a ':' deliminating hours and minutes")
	}

	hour, err := strconv.Atoi(hour_and_minute[0])
	if err != nil {
		return -1, -1, errors.New("hour is an invalid number")
	}

	if hour > 24 {
		return -1, -1, errors.New("hour > 24")
	}

	minute, err := strconv.Atoi(hour_and_minute[1])
	if err != nil {
		return -1, -1, errors.New("minute is an invalid number")
	}

	if minute >= 60 {
		return -1, -1, errors.New("minute >= 60")
	}

	return hour, minute, nil
}

func ParseCheckMessage(message string) (time.Time, error) {
	nowRegex := regexp.MustCompile("now")
	// this will match 1+ sequence of characters, an optional ':' and
	// additional 1+sequence of characters
	// we are intentionally making this loose in order to give precise feedback
	hhmmRegex := regexp.MustCompile("([0-9a-zA-Z]*)(:)?([0-9a-zA-Z]*)")
	currentTime := time.Now()

	if nowRegex.FindString(message) == "now" {
		return time.Now(), nil
	}

	date := hhmmRegex.FindString(message)
	if date == "" {
		return time.Time{}, fmt.Errorf("invalid time supplied: unsupported format")
	}

	hour, minute, err := parseTime(date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time supplied: %s", err)
	}
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), hour, minute, 0, 0, time.UTC), nil

}
