package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var playCommand = &discordgo.ApplicationCommand{
    Name:        "play",
    Description: "Play audio from a YouTube link",
    Options: []*discordgo.ApplicationCommandOption{
        {
            Type:        discordgo.ApplicationCommandOptionString,
            Name:        "url",
            Description: "YouTube video URL to play",
            Required:    true,
        },
    },
}

func RegisterHandlers(session *discordgo.Session) {
	session.AddHandler(onInteractionCreate)

	_, err := session.ApplicationCommandCreate(
		session.State.User.ID,
		"",
		playCommand,
	)
	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
}