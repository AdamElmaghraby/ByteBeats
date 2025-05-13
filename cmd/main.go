package main

import (
	"log"
	"os"
	"github.com/AdamElmaghraby/bytebeats/internal/discord"
	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == ""{
		log.Fatalf("Missing BOT_TOKEN env variable")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	discordgo.RegisterHandler(dg)


}