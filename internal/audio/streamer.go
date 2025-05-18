package audio

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jogramming/dca"
)

func StreamYouTube(ctx context.Context, url string, vc *discordgo.VoiceConnection) error {
	
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"

	encodingSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		return fmt.Errorf("dca encode error: %w", err)
	}
	defer encodingSession.Cleanup()

	vc.Speaking(true)
	defer vc.Speaking(false)

	done := make(chan error)
	dca.NewStream(encodingSession, vc, done)

	select{
	case err := <-done:
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("Stream error: %w", err)
		}
	case <-ctx.Done():
		return fmt.Errorf("playback timed out after %v", 10*time.Minute)
	}


	return nil
}