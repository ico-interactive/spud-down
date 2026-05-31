package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {

	// load env
	env, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("error: could not load .env file")
	}
	port := env["BOT_PORT"]
	token := env["DISCORD_TOKEN"]
	_ = []string{port}

	// setup discord bot session
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error: could not create discord session: ", err)
	}

	// destruct discord session
	defer discord.Close()

}
