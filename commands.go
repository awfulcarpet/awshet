package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	timeOption = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "time",
		Description: "a value, either 'now' or in the format 'xx:yy' to log",
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

func sendStringResponse(mesg string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: mesg,
			},
		},
	)
	if err != nil {
		log.Println("Unable to send response")
	}
}

func checkin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	time := i.ApplicationCommandData().GetOption("time").StringValue()
	// TODO: Null checks for non-guild servers
	log.Printf("%s invoked /checkin", i.Member.User.Username)

	checkinTime, err := parseCheckMessage(time)
	if err != nil {
		sendStringResponse(fmt.Sprintf(":x: %s", err.Error()), s, i)
		return
	}

	successMessage := fmt.Sprintf(":white_check_mark: Checked in %s at %02d:%02d on %02d/%02d\n", i.Member.User.Username,
		checkinTime.Hour(), checkinTime.Minute(), checkinTime.Month(), checkinTime.Day())

	sendStringResponse(successMessage, s, i)
}
