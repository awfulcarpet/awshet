package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	timeOption = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "time",
		Description: "a date value in the form mm/dd or now to record time",
		Required:    true,
	}
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "checkin",
			Description: "Checks students in and logs time and day",
			Options: []*discordgo.ApplicationCommandOption{
				timeOption,
			},
		},
		{
			Name:        "checkout",
			Description: "Checks students out and logs time and day",
			Options: []*discordgo.ApplicationCommandOption{
				timeOption,
			},
		},
	}
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
)

func registerCommands(dg *discordgo.Session) {
	log.Println("Creating all commands")
	// for i, v := range commands {
	// 	cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
	// 	if err != nil {
	// 		log.Panicf("Cannot create '%v' command: %v", v.Name, err)
	// 	}
	// 	registeredCommands[i] = cmd
	// }

	dg.AddHandler(slashCommandHandler)
}

func removeCommands(dg *discordgo.Session) {
	log.Println("Removing all Commands")

	for _, v := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func slashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "checkin":
		checkin(s, i)
	case "checkout":
		checkin(s, i)
	}
}

func checkin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	for _, v := range options {
		log.Println(v)
	}
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello world!",
			},
		},
	)
	if err != nil {
		log.Println("Unable to send response")
	}
}
