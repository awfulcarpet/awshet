package main

import (
	"github.com/bwmarrin/discordgo"
	"log"

	"time"
)

type check struct {
	name      string
	discordID string
	messageID string
	date      time.Time
	checkType string
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Println("message:", m.Message.ID, m.Author.Username, "sent:", m.Message.Content)
	handleCheck(s, m)
}

func handleCheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	checkType, date, err := parseCheckMessage(m.Message.Content)
	check := check{
		name:      m.Author.Username,
		discordID: m.Author.ID,
		messageID: m.Message.ID,
		date:      date,
		checkType: checkType,
	}
	if err != nil {
		s.MessageReactionAdd(m.ChannelID, m.Message.ID, "❌")
		log.Println("ERROR: unable to check in", check.name, "in message", check.messageID, err)
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.Message.ID, "✅")
	log.Println("checked", check.checkType, check.name, "on", check.date.Unix(), "in message", check.messageID)
}
