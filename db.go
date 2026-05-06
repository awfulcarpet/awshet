package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gocarina/gocsv"
)

type checkType string

type checkMessage struct {
	DiscordID string    `csv:"Discord ID"`
	Username  string    `csv:"Username"`
	Name      string    `csv:"Name"`
	Timestamp int64     `csv:"Timestamp"`
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

	return nil
}

func writeLog(msg checkMessage) error {
	f, err := os.OpenFile(CheckLogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("unable to open check log file: %s", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("unable to stat log file: %s", err)
	}

	// the file was created for the first time, need to add header for gocsv to
	// properly read the csv file
	if info.Size() == 0 {
		fmt.Fprintln(f, "Discord ID,Username,Name,Timestamp,Check Type")
	}

	_, err = fmt.Fprintf(f, "%s,%s,%s,%d,%s\n", msg.DiscordID, msg.Username, msg.Name,
		msg.Timestamp, msg.CheckType)
	if err != nil {
		return fmt.Errorf("unable to write mesg to check log file: %s", err)
	}

	return nil
}

func readLog() ([]*checkMessage, error) {
	f, err := os.OpenFile(CheckLogFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open users file for reading: %s", err)
	}
	defer f.Close()

	var logs []*checkMessage
	err = gocsv.UnmarshalFile(f, &logs)
	if err == io.EOF || err == gocsv.ErrEmptyCSVFile {
		return logs, nil
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse users file (%s): %s", f.Name(), err)
	}

	return logs, nil
}

func calculateTime(discordID string) (int, float32, error) {
	logs, err := readLog()
	if err != nil {
		return 0, 0, err
	}

	var hours float32 = 0.0

	var curTime int64 = 0.0

	for _, l := range logs {
		if l.DiscordID != discordID {
			continue
		}

		if l.CheckType == "in" {
			curTime = l.Timestamp
		}

		if l.CheckType == "out" {
			hours += float32(l.Timestamp-curTime) / 3600.0
			curTime = 0.0
		}
	}

	return 0, hours, nil
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
