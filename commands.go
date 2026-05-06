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
		{
			Name:        "time",
			Description: "lists currently logged time for user",
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
		check("in", s, i)
	case "checkout":
		check("out", s, i)
	case "time":
		getTime(s, i)
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
	}

	log.Printf("sent '%s' to discord\n", mesg)
}

func check(checkType checkType, s *discordgo.Session, i *discordgo.InteractionCreate) {
	time := i.ApplicationCommandData().GetOption("time").StringValue()

	if i.Member == nil {
		sendStringResponse(":x: awshet does not answer DMs", s, i)
		return
	}

	log.Printf("%s invoked /check%s %s", i.Member.User.Username, checkType, time)

	checkinTime, err := parseCheckMessage(time)
	if err != nil {
		sendStringResponse(fmt.Sprintf(":x: %s", err.Error()), s, i)
		return
	}

	msg := checkMessage{
		Username:  i.Member.User.Username,
		Name:      i.Member.Nick,
		DiscordID: i.Member.User.ID,
		Timestamp: checkinTime.Unix(),
		CheckType: checkType,
	}

	err = updateDB(msg)

	if err != nil {
		sendStringResponse(fmt.Sprintf(":x: %s", err.Error()), s, i)
		return
	}

	sendStringResponse(
		fmt.Sprintf(":white_check_mark: Checked %s %s (%s) at %02d:%02d on %02d/%02d",
			checkType, msg.Name, msg.Username, checkinTime.Hour(),
			checkinTime.Minute(), checkinTime.Month(), checkinTime.Day()),
		s, i)
}

func getTime(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		sendStringResponse(":x: awshet does not answer DMs", s, i)
		return
	}

	log.Printf("%s invoked /time", i.Member.User.Username)

	days, hours, err := calculateTime(i.Member.User.ID)

	if err != nil {
		sendStringResponse(fmt.Sprintf(":x: %s", err), s, i)
		return
	}

	sendStringResponse(
		fmt.Sprintf("%s has logged %.2f hours over the course of %d days",
			i.Member.Nick, hours, days), s, i)
}
