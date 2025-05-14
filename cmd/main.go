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

	err := godotenv.Load()
	if err != nil {
        log.Println("No .env file found")
    }

    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("Missing BOT_TOKEN")
    }

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()
	discord.RegisterHandlers(dg)
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