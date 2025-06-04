package discord

import (
	"log"

	"github.com/AdamElmaghraby/ByteBeats/internal/audio"
	"github.com/bwmarrin/discordgo"
)

// RegisterHandler registers the slash command handlers
func RegisterHandler(dg *discordgo.Session) {
	dg.AddHandler(interactionCreate)
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Only handle /play commands
	if i.Type != discordgo.InteractionApplicationCommand || i.ApplicationCommandData().Name != "play" {
		return
	}

	// Extract URL
	url := i.ApplicationCommandData().Options[0].StringValue()

	// ACK to avoid timeout
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		log.Println("ACK error:", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Failed to acknowledge command.",
		})
		return
	}

	// Find user's voice channel
	vs := findUserVoiceState(s, i.GuildID, i.Member.User.ID)
	if vs == nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ You must be in a voice channel to use /play.",
		})
		return
	}

	// Join voice
	vc, err := s.ChannelVoiceJoin(i.GuildID, vs.ChannelID, false, true)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Could not join your voice channel.",
		})
		log.Println("Voice join error:", err)
		return
	}

	// Kick off audio playback; PlayAudio will send its own follow-ups on success or error
	go audio.PlayAudio(s, vc, url, i)
}

// findUserVoiceState returns the voice state for a user in a guild
func findUserVoiceState(s *discordgo.Session, guildID, userID string) *discordgo.VoiceState {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil
	}
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs
		}
	}
	return nil
}
