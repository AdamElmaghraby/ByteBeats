package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	if data.Name != "play" {
		return
	}

	url := data.Options[0].StringValue()

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
	
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Loading your beat...",
		},
	})

	go func() {
		vc, err := s.ChannelVoiceJoin(i.GuildID, vs.ChannelID, false, true)
		if err != nil {
			msg := fmt.Sprintf("Failed to join voice channel: %v", err)
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{Content: msg})
			return
		}
		defer vc.Disconnect()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		err = audio.StreamYouTube(ctx, url, vc)
		if err != nil {
			msg := fmt.Sprintf("Playback error: %v", err)
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{Content: msg})
			return
		}

		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: "Finished playing your beats!",
		})


	}()
	

}