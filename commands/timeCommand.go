package commands

import (
	"awshet/db"
	"cmp"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
)

func calculateTime(discordID string) (int, float32, error) {
	logs, err := db.ReadLog()
	if err != nil {
		return 0, 0, err
	}

	slices.SortStableFunc(logs, func(a, b *db.CheckMessage) int {
		return cmp.Compare(a.Timestamp, b.Timestamp)
	})

	var hours float32 = 0.0

	days := make(map[time.Time]bool)

	var curTime int64 = 0.0

	for _, l := range logs {
		if l.DiscordID != discordID {
			continue
		}

		year, month, day := time.Unix(l.Timestamp, 0).Local().Date()
		k := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		days[k] = true

		if l.CheckType == "in" {
			curTime = l.Timestamp
		}

		if l.CheckType == "out" {
			hours += float32(l.Timestamp-curTime) / 3600.0
			curTime = 0.0
		}
	}

	return len(days), hours, nil
}

func timeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
