package discord

import (
	"context"

	"github.com/AdamElmaghraby/ByteBeats/internal/audio"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	if data.Name != "play" {
		return
	}

	url := data.Options[0].StringValue()
	
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	vs, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
    if err != nil || vs == nil || vs.ChannelID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You need to first be in a voice channel to hear your beats!",
			},
		})
		return
	} 

	vc, err := s.ChannelVoiceJoin(guildID, vs.ChannelID, false, true)
    if err != nil {
        s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
            Content: "Failed to join voice channel.",
        })
        return
    }
	
	go func() {
        err := audio.StreamURL(context.Background(), url, vc)
        if err != nil {
            s.ChannelMessageSend(vs.ChannelID, "Error streaming audio: "+err.Error())
        }
        vc.Disconnect()
    }()

    s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
        Content: "▶️ Streaming now!",
    })

}