package main

import (
	"fmt"
	"log"
	"os"

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

	log.Println("message:", m.Message.ID, m.Author.Username, "sent", m.Message.Content)
	s.MessageReactionAdd(m.ChannelID, m.Message.ID, "✅")

	// Respond to a simple command "!ping"
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
