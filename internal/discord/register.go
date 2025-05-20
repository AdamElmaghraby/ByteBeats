package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)



func RegisterCommand(session *discordgo.Session) error {
	_, err := session.ApplicationCommandCreate(s.State.User.ID, "",	&discordgo.ApplicationCommand{
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
	
	return err

}