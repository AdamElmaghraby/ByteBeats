package audio

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/bwmarrin/discordgo"
)

func StreamYouTube(ctx context.Context, url string, vc *discordgo.VoiceConnection) error {
	ytdlp := exec.CommandContext(ctx, "yt-dlp", "-f", "bestaudio", "-o", "-", url)
	
	ffmpeg := exec.CommandContext(ctx,
        "ffmpeg", "-i", "pipe:0",
        "-analyzeduration", "0", "-loglevel", "0",
        "-ac", "2", "-ar", "48000",
        "-f", "s16le", "pipe:1",
    )

	r, w := io.Pipe()
    ytdlp.Stdout = w
    ffmpeg.Stdin = r

	pcmOut, err := ffmpeg.StdoutPipe()
    if err != nil {
        return err
    }

    if err := ytdlp.Start(); err != nil {
        return err
    }
    if err := ffmpeg.Start(); err != nil {
        return err
    }

	vc.Speaking(true)
	defer vc.Speaking(false)

	vc.Speaking(true)
    defer vc.Speaking(false)

    frameSize := 1920
    pcmBuf := make([]byte, frameSize*2*2) 
    opusBuf := make([]byte, 4000)

	for {
        n, err := pcmOut.Read(pcmBuf)
        if err != nil {
            break
        }
        opusFrame, err := vc.OpusEncode(pcmBuf[:n], opusBuf)
        if err != nil {
            return err
        }
        vc.OpusSend <- opusBuf[:opusFrame]
    }
    return nil
}