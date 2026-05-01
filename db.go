package main

import (
	"fmt"
	"log"
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

func writeLog(msg checkMessage) {
	f, err := os.OpenFile(CheckLogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("ERR: unable to open %s: %s\n", CheckLogFileName, err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s,%s,%d,%s\n", msg.discordID, msg.username,
		msg.time.Unix(), msg.checkType)
	if err != nil {
		log.Printf("ERR: Unable to write mesg to %s: %s\n", CheckLogFileName, err)
	}

	userLogs := readUserLog()
	writeNewUserLogs(userLogs)
}

func writeNewUserLogs(userLogs []*userLog) {
	writeUserLog(userLogs)
}

func readUserLog() []*userLog {
	f, err := os.OpenFile(UsersLogfileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Printf("ERR: Unable to open file %s: %s\n", UsersLogfileName, err)
	}
	defer f.Close()

	var userLogs []*userLog
	if err = gocsv.UnmarshalFile(f, &userLogs); err != nil {
		log.Printf("ERR: Unable to parse file %s: %s\n", UsersLogfileName, err)
	}

	return userLogs
}

func writeUserLog(userLogs []*userLog) {
	f, err := os.OpenFile(UsersLogfileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	// Check for errors when opening or creating the file. If there's an error, panic.
	if err != nil {
		log.Printf("ERR: Unable to open file %s: %s\n", UsersLogfileName, err)
	}
	defer f.Close()

	// Marshal the userLogs into the CSV format and write them to the result.csv file
	if err = gocsv.MarshalFile(&userLogs, f); err != nil {
		log.Printf("ERR: Unable to parse file %s: %s\n", UsersLogfileName, err)
	}
}
