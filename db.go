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
	Timestamp int64     `csv:"Timestamp"`
	CheckType checkType `csv:"Check Type"`
}

var (
	CheckLogFileName = "checklog.csv"
	UsersLogfileName = "users.csv"
)

func WriteLog(msg checkMessage) error {
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

func ReadLog() ([]*checkMessage, error) {
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

func CalculateTime(discordID string) (int, float32, error) {
	logs, err := ReadLog()
	if err != nil {
		return 0, 0, err
	}

	var hours float32 = 0.0

	days := make(map[time.Time]bool)

	var curTime int64 = 0.0

	for _, l := range logs {
		if l.DiscordID != discordID {
			continue
		}

		year, month, day := time.Unix(l.Timestamp, 0).Date()
		k := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		days[k] = true

		if l.CheckType == "in" {
			curTime = l.Timestamp
		}

		if l.CheckType == "out" {
			hours += float32(l.Timestamp-curTime) / 3600.0
			curTime = 0.0
		}
	}

	return len(days), hours, nil
}
