package audio

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

func StreamYouTube(ctx context.Context, url string, vc *discordgo.VoiceConnection) error {

	ytdlp := exec.CommandContext(ctx, "yt-dlp", "-f", "bestaudio", "-o", "-", url)

	ffmpeg := exec.CommandContext(ctx,
        "ffmpeg",
        "-i", "pipe:0",       
        "-analyzeduration", "0",
        "-loglevel", "0",
        "-ac", "2",           
        "-ar", "48000",       
        "-f", "s16le",        
        "pipe:1",             
    )

	r, w := io.Pipe()
	ytdlp.Stdout = w
	ffmpeg.Stdin = r

	pcmOut, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}

	err = ytdlp.Start()
	if err != nil {
		return err
	}

	err = ffmpeg.Start()
	if err != nil {
		return err
	}

	vc.Speaking(true)
	defer vc.Speaking(false)

	opusEncoder, err := gopus.NewEncoder(48000, 2, gopus.Audio)
	if err != nil {
    	return fmt.Errorf("could not create gopus encoder: %w", err)
	}

	const frameSize = 1920
	pcmBuf := make([]byte, frameSize)
	opusBuf := make([]byte, 4000)

	for {
		n, err := io.ReadFull(pcmOut, pcmBuf)
		if err != nil {
			break
		}

		frameSize, err := opusEncoder.Encode(pcmBuf[:n], opusBuf)
			if err != nil {
    			return fmt.Errorf("opus encode error: %w", err)
			}

		vc.OpusSend <- opusBuf[:frameSize]

	}

	w.Close()
	ytdlp.Wait()
	ffmpeg.Wait()

	return nil
}