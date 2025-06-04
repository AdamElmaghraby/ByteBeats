package audio

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// PlayAudio is a minimal implementation for Discord streaming.
// s: Discord session
// vc: joined VoiceConnection
// query: user-supplied YouTube URL
// i: original InteractionCreate, so we can send follow-up messages
func PlayAudio(s *discordgo.Session, vc *discordgo.VoiceConnection, query string, i *discordgo.InteractionCreate) {
	os.Setenv("DCA_FFMPEG_OPTS", "-hide_banner -loglevel debug -map 0:a")
	log.Printf("[PlayAudio] Start for query: %s", query)

	// 1) Validate URL
	u, err := url.Parse(query)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		log.Printf("[PlayAudio] URL parse error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Invalid URL.",
		})
		vc.Disconnect()
		return
	}

	host := strings.ToLower(u.Hostname())
	if !strings.Contains(host, "youtube.com") && !strings.Contains(host, "youtu.be") {
		log.Printf("[PlayAudio] Invalid host: %s", host)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Not a YouTube URL.",
		})
		vc.Disconnect()
		return
	}

	// 2) Extract video ID
	var videoID string
	if strings.Contains(host, "youtu.be") {
		videoID = strings.TrimPrefix(u.Path, "/")
	} else {
		videoID = u.Query().Get("v")
	}
	if videoID == "" {
		log.Println("[PlayAudio] Missing video ID")
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Could not determine video ID.",
		})
		vc.Disconnect()
		return
	}
	log.Printf("[PlayAudio] Video ID: %s", videoID)

	// 3) Download WebM/Opus via yt-dlp
	tmpWebM, err := os.CreateTemp("", "audio-*.webm")
	if err != nil {
		log.Printf("[PlayAudio] Temp file error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Could not create temp file.",
		})
		vc.Disconnect()
		return
	}
    tmpWebM.Close()
	defer func() {
		/*tmpWebM.Close()*/
		os.Remove(tmpWebM.Name())
	}()

	watchURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	downloadCmd := exec.Command(
		".//yt-dlp",
		"--rm-cache-dir",
		"--no-cache-dir",
		"--force-overwrites",
		"-f", "bestaudio[ext=webm]",
		"-o", tmpWebM.Name(),
		watchURL,
	)
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	log.Printf("[PlayAudio] Running download: %v", downloadCmd.Args)
	if err := downloadCmd.Run(); err != nil {
		log.Printf("[PlayAudio] yt-dlp error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("❌ Download error: %v", err),
		})
		vc.Disconnect()
		return
	}

	infoWebM, err := os.Stat(tmpWebM.Name())
	if err != nil || infoWebM.Size() == 0 {
		log.Printf("[PlayAudio] Download failed or file empty: %v, size: %d", err, infoWebM.Size())
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Download resulted in empty file.",
		})
		vc.Disconnect()
		return
	}
	log.Printf("[PlayAudio] WebM download complete, size: %d bytes", infoWebM.Size())

	// 4) Convert WebM → pure Ogg/Opus via ffmpeg
	tmpOgg := tmpWebM.Name() + ".opus"
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", tmpWebM.Name(),
		"-c:a", "libopus",
		"-f", "opus",
		tmpOgg,
	)
	ffmpegCmd.Stdout = os.Stdout
	ffmpegCmd.Stderr = os.Stderr
	log.Printf("[PlayAudio] Running ffmpeg re-encode: %v", ffmpegCmd.Args)
	if err := ffmpegCmd.Run(); err != nil {
		log.Printf("[PlayAudio] FFmpeg re-encode error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("❌ ffmpeg error: %v", err),
		})
		vc.Disconnect()
		return
	}
	defer os.Remove(tmpOgg)

	infoOgg, err := os.Stat(tmpOgg)
	if err != nil || infoOgg.Size() == 0 {
		log.Printf("[PlayAudio] ffmpeg produced empty file: %v, size %d", err, infoOgg.Size())
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ ffmpeg re-encode resulted in empty file.",
		})
		vc.Disconnect()
		return
	}
	log.Printf("[PlayAudio] ffmpeg re-encode complete, size: %d bytes", infoOgg.Size())

	// 5) Notify user that playback is starting
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "▶️ Now playing!",
	})

	// 6) Encode & stream with DCA
	options := dca.StdEncodeOptions
	options.BufferedFrames   = 100
	options.FrameDuration    = 20
	options.CompressionLevel = 5
	options.Bitrate          = 96

	log.Println("[PlayAudio] Encoding Ogg/Opus file for streaming…")
	pathToFFMPEG, lookErr := exec.LookPath("ffmpeg")
	log.Printf("[PlayAudio] DCA will use ffmpeg from: %s (lookErr=%v)", pathToFFMPEG, lookErr)

	encodeSession, err := dca.EncodeFile(tmpOgg, options)
	if err != nil {
		log.Printf("[PlayAudio] Failed to create encoding session for '%s': %v", tmpOgg, err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("❌ Encoding error: %v", err),
		})
		vc.Disconnect()
		return
	}
	defer func() {
		encodeSession.Cleanup()
		log.Println("[PlayAudio] Cleared EncodeSession")
	}()

	// Give DCA/FFmpeg a brief moment to buffer before speaking
	time.Sleep(500 * time.Millisecond)

	log.Println("[PlayAudio] Starting DCA.NewStream…")
	doneChan := make(chan error)
	stream := dca.NewStream(encodeSession, vc, doneChan)
	if stream == nil {
		log.Println("[PlayAudio] DCA.NewStream returned nil—stream init failed")
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Stream initialization failed.",
		})
		vc.Disconnect()
		return
	}

	vc.Speaking(true)
	log.Println("[PlayAudio] vc.Speaking(true) called")

	// 7) Wait for streaming to finish (EOF) or for an error
	if err := <-doneChan; err != nil && err != io.EOF {
		log.Printf("[PlayAudio] Streaming error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("❌ Stream error: %v", err),
		})
	} else {
		log.Println("[PlayAudio] Stream finished (EOF or no error)")
	}

	vc.Speaking(false)
	log.Println("[PlayAudio] vc.Speaking(false) called")
	vc.Disconnect()
	log.Println("[PlayAudio] Playback finished and disconnected")
}
