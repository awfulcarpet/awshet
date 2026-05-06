package commands

import (
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
		{
			Name:        "time",
			Description: "lists currently logged time for user",
		},
	}
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
)

func RegisterCommands(dg *discordgo.Session) {
	log.Println("Creating all commands")
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	dg.AddHandler(slashCommandHandler)
}

func RemoveCommands(dg *discordgo.Session) {
	log.Println("Removing all Commands")

	for _, v := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, v.GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func slashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "checkin":
		checkCommand("in", s, i)
	case "checkout":
		checkCommand("out", s, i)
	case "time":
		timeCommand(s, i)
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
		log.Println("ERR: Unable to send response")
		return
	}

	log.Printf("sent '%s' to discord\n", mesg)
}
