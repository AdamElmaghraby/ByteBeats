package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)



func RegisterCommand(dg *discordgo.Session) error {
	_, err := dg.ApplicationCommandCreate(dg.State.User.ID, "",	&discordgo.ApplicationCommand{
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
    })
	if err != nil {
        log.Fatalf("Cannot create slash cmd: %v", err)
    }

	return err

}