package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Get the bot token from environment variable (recommended)
	Token := os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		fmt.Println("No token provided. Set the DISCORD_TOKEN environment variable.")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate function as a callback for the MessageCreate event.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	fmt.Println("Bot is now running")
	select {}
}

// This function will be called every time a new message is created in the Discord server.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Println("message:", m.Message.ID, m.Author.Username, "sent:", m.Message.Content)
	handleCheck(s, m)
}

type check struct {
	name      string
	discordID string
	messageID string
	date      uint64
	checkType string
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
	log.Println("checked in", check.name, "in message", check.messageID)
}

func parseCheckMessage(message string) (string, uint64, error) {
	r, _ := regexp.Compile("(in|out)")

	match := r.MatchString(message)
	if !match {
		return "", 0, errors.New("invalid command supplied")
	}

	return "", 0, nil
}
