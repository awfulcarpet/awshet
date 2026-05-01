package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	logOutput = flag.String("f", "/dev/stderr", "the location to write the logs to")
)

func main() {
	flag.Parse()

	f, err := os.OpenFile(*logOutput, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error log file: %v", err)
	}
	log.SetOutput(f)
	log.Println("logging to", *logOutput)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Token := os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		fmt.Println("No token provided. Set the DISCORD_TOKEN environment variable.")
		return
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	registerCommands(dg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Shutting Down")
	// removeCommands(dg)
}
