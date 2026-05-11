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

func formatLogs(discordID string) (string, error) {
	logs, err := db.ReadLog()
	if err != nil {
		return "", err
	}

	slices.SortStableFunc(logs, func(a, b *db.CheckMessage) int {
		return cmp.Compare(a.Timestamp, b.Timestamp)
	})

	type pair struct {
		in  int64
		out int64
	}

	var output string = ""
	days := make(map[time.Time]pair)

	for _, l := range logs {
		if l.DiscordID != discordID {
			continue
		}

		year, month, day := time.Unix(l.Timestamp, 0).Local().Date()
		k := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

		if l.CheckType == "in" {
			days[k] = pair{
				in:  l.Timestamp,
				out: days[k].out,
			}
		}

		if l.CheckType == "out" {
			days[k] = pair{
				in:  days[k].in,
				out: l.Timestamp,
			}
			in_hours, in_minutes, _ := time.Unix(days[k].in, 0).Clock()
			out_hours, out_minutes, _ := time.Unix(days[k].out, 0).Clock()
			output = fmt.Sprintf("%s%02d/%02d/%02d: in: %02d:%02d | out: %02d:%02d\n", output, year, month, day, in_hours, in_minutes, out_hours, out_minutes)
		}

	}

	return output, nil
}

func listLogCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		sendStringResponse(":x: awshet does not answer DMs", s, i)
		return
	}

	log.Printf("%s invoked /listlog", i.Member.User.Username)

	output, err := formatLogs(i.Member.User.ID)

	if err != nil {
		sendStringResponse(fmt.Sprintf(":x: %s", err), s, i)
		return
	}

	sendStringResponse(output, s, i)
}
