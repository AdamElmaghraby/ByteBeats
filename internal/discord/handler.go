package discord

import (
	"log"

	"github.com/AdamElmaghraby/ByteBeats/internal/audio"
	"github.com/bwmarrin/discordgo"
)

func RegisterHandler(dg *discordgo.Session) {
    dg.AddHandler(interactionCreate)
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.Type != discordgo.InteractionApplicationCommand {
        return
    }

    data := i.ApplicationCommandData()
    if data.Name != "play" {
        return
    }

    url := data.Options[0].StringValue()

    // 1) ACK to avoid timeout
    if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
    }); err != nil {
        log.Println("ACK error:", err)
        return
    }

    // 2) Find & join user voice channel
    vs := findUserVoiceState(s, i.GuildID, i.Member.User.ID)
    if vs == nil {
        s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
            Content: "You must be in a voice channel!",
        })
        return
    }
    vc, err := s.ChannelVoiceJoin(i.GuildID, vs.ChannelID, false, true)
    if err != nil {
        s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
            Content: "❌ Failed to join voice channel.",
        })
        return
    }

    // 3) Stream audio
    go audio.PlayAudio(s, vc, url)

    // 4) Confirm
    s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
        Content: "▶️ Now playing!",
    })
}

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