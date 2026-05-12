package commands

import (
	"awshet/db"
	"cmp"
	"fmt"
	"log"
	"math"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	// this is the default checkin/checkout times, only fill in hour + minutes
	// as the rest are overridden
	defaultCheckInTime  = time.Date(0, 0, 0, 16, 0, 0, 0, time.Local)
	defaultCheckOutTime = time.Date(0, 0, 0, 18, 0, 0, 0, time.Local)
)

type pair = struct {
	in  int64
	out int64
}

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

	var check pair

	for _, l := range logs {
		if l.DiscordID != discordID {
			continue
		}
		switch l.CheckType {
		case "in":
			check.in = l.Timestamp
		case "out":
			check.out = l.Timestamp

			year, month, day := time.Unix(l.Timestamp, 0).Local().Date()
			k := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
			days[k] = true

			if check.in == 0 {
				check.in = defaultCheckInTime.AddDate(year, int(month), day).Unix()
			}

			timeDiff := float32(check.out-check.in) / 3600.0
			if timeDiff < 0 {
				return 0, 0, fmt.Errorf("time diff calculated on %d/%d/%d (%d) is negative (%.2f)", year, month, day, l.Timestamp, timeDiff)
			}
			hours += timeDiff

			check = pair{}
		default:
			return 0, 0, fmt.Errorf("checktype for log is somehow incorrect")
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

	hour, min := math.Modf(float64(hours))

	sendStringResponse(
		fmt.Sprintf("%s has logged %2d hours and %2d min over the course of %d days",
			i.Member.Nick, int(hour), int(min*60), days), s, i)
}
