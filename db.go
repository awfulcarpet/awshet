package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type checkType string

type checkMessage struct {
	username  string    `csv:"Username"`
	discordID string    `csv:"Discord ID"`
	time      time.Time `csv:"Timestamp"`
	checkType checkType `csv:"Check Type"`
}

type userLog struct {
	Username   string  `csv:"Username"`
	DiscordID  string  `csv:"Discord ID"`
	TotalDays  int     `csv:"Total Days"`
	TotalHours float32 `csv:"Total Hours"`
}

var (
	CheckLogFileName = "checklog.csv"
	UsersLogfileName = "users.csv"
)

func updateDB(msg checkMessage) error {
	err := writeLog(msg)
	if err != nil {
		return err
	}

	userLogs, err := readUserLog()
	if err != nil {
		return err
	}

	err = writeNewUserLogs(userLogs)
	if err != nil {
		return err
	}

	return nil
}

func writeLog(msg checkMessage) error {
	f, err := os.OpenFile(CheckLogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("unable to open check log file: %s\n", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s,%s,%d,%s\n", msg.discordID, msg.username,
		msg.time.Unix(), msg.checkType)
	if err != nil {
		return fmt.Errorf("unable to write mesg to check log file: %s\n", err)
	}

	return nil
}

// TODO (may be unneeded) write logic for parsing total days and hours and write
func writeNewUserLogs(userLogs []*userLog) error {
	return nil
	// return writeUserLog(userLogs)
}

func readUserLog() ([]*userLog, error) {
	f, err := os.OpenFile(UsersLogfileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("unable to open users file for reading: %s\n", err)
	}
	defer f.Close()

	var userLogs []*userLog
	if err = gocsv.UnmarshalFile(f, &userLogs); err != nil {
		return nil, fmt.Errorf("unable to parse users file (%s): %s\n", UsersLogfileName, err)
	}

	return userLogs, nil
}

func writeUserLog(userLogs []*userLog) error {
	f, err := os.OpenFile(UsersLogfileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	// Check for errors when opening or creating the file. If there's an error, panic.
	if err != nil {
		return fmt.Errorf("unable to open users file for writing: %s\n", err)
	}
	defer f.Close()

	// Marshal the userLogs into the CSV format and write them to the result.csv file
	if err = gocsv.MarshalFile(&userLogs, f); err != nil {
		return fmt.Errorf("unable to write to users file: %s\n", err)
	}
	return nil
}
