package audio

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jogramming/dca"
)

func StreamYouTube(ctx context.Context, url string, vc *discordgo.VoiceConnection) error {
	cmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", "-o", "temp_audio.mp3", url)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("youtube-dl error: %w", err)
	}
	defer os.Remove("temp_audio.mp3")

	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"

	encodingSession, err := dca.EncodeFile("temp_audio.mp3", opts)
	if err != nil {
		return fmt.Errorf("dca encode error: %w", err)
	}
	defer encodingSession.Cleanup()

	vc.Speaking(true)
	defer vc.Speaking(false)

	done := make(chan error)
	dca.NewStream(encodingSession, vc, done)

	select {
	case err := <-done:
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("stream error: %w", err)
		}
	case <-ctx.Done():
		return fmt.Errorf("playback timed out after %v", 10*time.Minute)
	}

	return nil
}