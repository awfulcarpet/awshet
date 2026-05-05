package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type checkType string

type checkMessage struct {
	DiscordID string    `csv:"Discord ID"`
	Username  string    `csv:"Username"`
	Name      string    `csv:"Name"`
	Timestamp time.Time `csv:"Timestamp"`
	CheckType checkType `csv:"Check Type"`
}

type userLog struct {
	DiscordID  string  `csv:"Discord ID"`
	Username   string  `csv:"Username"`
	Name       string  `csv:"Name"`
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
		return fmt.Errorf("unable to open check log file: %s", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s,%s,%s,%d,%s\n", msg.DiscordID, msg.Username, msg.Name,
		msg.Timestamp.Unix(), msg.CheckType)
	if err != nil {
		return fmt.Errorf("unable to write mesg to check log file: %s", err)
	}

	return nil
}

func readUserLog() ([]*userLog, error) {
	f, err := os.OpenFile(UsersLogfileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open users file for reading: %s", err)
	}
	defer f.Close()

	var userLogs []*userLog
	err = gocsv.UnmarshalFile(f, &userLogs)
	if err == io.EOF || err == gocsv.ErrEmptyCSVFile {
		return userLogs, nil
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse users file (%s): %s", UsersLogfileName, err)
	}

	return userLogs, nil
}

// TODO (may be unneeded) write logic for parsing total days and hours and write
func writeNewUserLogs(userLogs []*userLog) error {
	return writeUserLog(userLogs)
}

func writeUserLog(userLogs []*userLog) error {
	f, err := os.OpenFile(UsersLogfileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	// Check for errors when opening or creating the file. If there's an error, panic.
	if err != nil {
		return fmt.Errorf("unable to open users file for writing: %s", err)
	}
	defer f.Close()

	// Marshal the userLogs into the CSV format and write them to the result.csv file
	if err = gocsv.MarshalFile(&userLogs, f); err != nil {
		return fmt.Errorf("unable to write to users file: %s", err)
	}
	return nil
}
