package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type checkType string

type checkMessage struct {
	username  string
	discordID string
	time      time.Time
	checkType checkType
}

var (
	LogFileName = "log.csv"
)

func writeLog(msg checkMessage) {
	f, err := os.OpenFile(LogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("ERR: unable to open %s: %s\n", LogFileName, err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s,%s,%d,%s\n", msg.discordID, msg.username,
		msg.time.Unix(), msg.checkType)
	if err != nil {
		log.Printf("ERR: Unable to write mesg to %s: %s\n", LogFileName, err)
	}
}
