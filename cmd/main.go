package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AdamElmaghraby/ByteBeats/internal/discord"
	"github.com/joho/godotenv"

	"github.com/bwmarrin/discordgo"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found; falling back to real env vars")
    }

    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("Missing BOT_TOKEN")
    }

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	discord.RegisterHandler(dg)

	if err := dg.Open(); err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()
	log.Println(`⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⡟⠷⣶⣤⠀⠀⠀⠀⠀
⠀⠀⠀⢸⠟⠛⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣼⡇⢀⣸⡇⠀⠀⠀⠀⠀
⠀⠀⠀⠘⣧⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⠿⠿⠻⣿⣿⡿⠀⠀⠀⠀⠀
⠀⠀⠀⣾⣿⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠈⠉⠀⠀⠀⠀⣿⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⢀⣀⣀⣀⣀⣀⣀⣀⣉⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣉⣀⣀⣀⣀⣀⣀⣀⡀⠀
⠀⢸⣿⣿⣉⣉⣉⣹⣿⣿⣏⣉⣉⣉⣉⣉⣉⣉⣉⣹⣿⣿⣏⣉⣉⣉⣿⣿⡇⠀
⠀⢸⣿⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⣿⡇⠀
⠀⢸⣿⡿⠋⣉⣭⣍⡙⢿⣿⡏⢉⣭⣭⣭⣭⡉⢹⣿⡿⢋⣩⣭⣉⠙⢿⣿⡇⠀
⠀⢸⡟⢠⣾⠟⠛⠻⣿⣆⠹⡇⢸⣿⣿⣿⣿⡇⢸⠏⣰⣿⠟⠛⠻⣷⡄⢻⡇⠀
⠀⢸⡇⢸⣿⡀⠛⢀⣿⡿⢀⣧⣤⣤⣤⣤⣤⣤⣼⡀⢿⣿⡀⠛⢀⣿⡇⢸⡇⠀
⠀⢸⣿⣄⠙⠿⣿⡿⠟⣡⣾⣿⡏⢹⡏⢹⡏⢹⣿⣷⣌⠻⢿⣿⠿⠋⣠⣿⡇⠀
⠀⠸⠿⠿⠷⠶⠶⠶⠾⠿⠿⠿⠿⠾⠿⠿⠷⠿⠿⠿⠿⠷⠶⠶⠶⠾⠿⠿⠇⠀
⠀⠀
ByteBeats is now up and running! Press CTRL-C to exit.`)

stop := make(chan os.Signal, 1)
signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
<-stop
	

}