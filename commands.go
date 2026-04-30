package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "checkin",
			Description: "Checks students in and logs time and day",
		},
		{
			Name:        "checkout",
			Description: "Checks students out and logs time and day",
		},
	}
)

func registerCommands(dg *discordgo.Session) {
	log.Println("Creating all commands")
	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
}

func removeCommands(dg *discordgo.Session) {
	log.Println("Removing all Commands")

	for _, v := range commands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
