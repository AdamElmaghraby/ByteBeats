package audio

import (
    "io"
    "log"
    "github.com/jonas747/dca"
    "github.com/bwmarrin/discordgo"
)

func PlayAudio(s *discordgo.Session, vc *discordgo.VoiceConnection, url string) {
    // 1) Encode + stream options
    opts := dca.StdEncodeOptions
    opts.RawOutput   = true
    opts.Bitrate     = 96
    opts.Application = "lowdelay"

    // 2) Create encoding session (yt-dlp + ffmpeg under the hood)
    encodeSession, err := dca.EncodeFile(url, opts)
    if err != nil {
        log.Println("Encode error:", err)
        return
    }
    defer encodeSession.Cleanup()

    // 3) Start the DCA stream into the voice connection
    done := make(chan error)
    _ = dca.NewStream(encodeSession, vc, done)

    vc.Speaking(true)
    defer vc.Speaking(false)

    // 4) Wait until stream finishes or errors
    if err := <-done; err != nil && err != io.EOF {
        log.Println("Stream error:", err)
    }

    // 5) after done, leave channel
    vc.Disconnect()
}
