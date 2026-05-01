package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type checkType string

type checkMessage struct {
	username  string
	discordID string
	time      time.Time
	checkType checkType
}

type userLog struct {
	username   string  `csv:"username"`
	discordID  string  `csv:"discordID"`
	totalDays  int     `csv:"totalDays"`
	totalHours float32 `csv:"totalHours"`
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

	fmt.Println(readUserLog())
}

func readUserLog() ([]userLog, error) {
	file, err := os.Open(UsersLogfileName)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)

	users := []userLog{}
	line, err := reader.Read()
	if err != nil {
		return users, err
	}

	if len(line) < 4 {
		return users, errors.New("too few columns in user log file")
	}

	totalDays, _ := strconv.Atoi(line[2])
	totalHours, _ := strconv.Atoi(line[3])

	users = append(users, userLog{
		username:   line[0],
		discordID:  line[1],
		totalDays:  totalDays,
		totalHours: float32(totalHours),
	})

	return users, nil
}
