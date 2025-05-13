package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	if data.Name != "play" {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Received `/play`! URL: " + data.Options[0].StringValue(),
		},
	})
	if err != nil {
		log.Printf("Failed to respond to interaction: %v", err)
	}
}