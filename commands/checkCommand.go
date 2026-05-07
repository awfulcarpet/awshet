package commands

import (
	"awshet/db"
	"awshet/parsing"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func checkCommand(checkType db.CheckType, s *discordgo.Session, i *discordgo.InteractionCreate) {
	time := i.ApplicationCommandData().GetOption("time").StringValue()

	if i.Member == nil {
		sendStringResponse(":x: awshet does not answer DMs", s, i)
		return
	}

	log.Printf("%s invoked /check%s %s", i.Member.User.Username, checkType, time)

	checkinTime, err := parsing.ParseCheckMessage(time)
	if err != nil {
		sendStringResponse(fmt.Sprintf(":x: %s", err.Error()), s, i)
		return
	}

	msg := db.CheckMessage{
		Username:  i.Member.User.Username,
		Name:      i.Member.Nick,
		DiscordID: i.Member.User.ID,
		Timestamp: checkinTime.Unix(),
		CheckType: checkType,
	}

	err = db.WriteMessageToLog(msg)

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
