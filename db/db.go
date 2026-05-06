package db

import (
	"fmt"
	"io"
	"os"

	"github.com/gocarina/gocsv"
)

type CheckType string

type CheckMessage struct {
	DiscordID string    `csv:"Discord ID"`
	Username  string    `csv:"Username"`
	Name      string    `csv:"Name"`
	Timestamp int64     `csv:"Timestamp"`
	CheckType CheckType `csv:"Check Type"`
}

var (
	CheckLogFileName = "checklog.csv"
	UsersLogfileName = "users.csv"
)

func WriteLog(msg CheckMessage) error {
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

func ReadLog() ([]*CheckMessage, error) {
	f, err := os.OpenFile(CheckLogFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open users file for reading: %s", err)
	}
	defer f.Close()

	var logs []*CheckMessage
	err = gocsv.UnmarshalFile(f, &logs)
	if err == io.EOF || err == gocsv.ErrEmptyCSVFile {
		return logs, nil
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse users file (%s): %s", f.Name(), err)
	}

	return logs, nil
}

