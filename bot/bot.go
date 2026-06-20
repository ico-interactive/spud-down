package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	fnafAPI        = "localhost:9638"
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
)

var session *discordgo.Session

func init() { flag.Parse() }

func init() {
	// load env
	env, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("error: could not load .env file")
	}

	token := env["DISCORD_TOKEN"]

	// setup session bot session
	session, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error: could not create discord session: ", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "merge",
			Description: "what follows?",
		},
		{
			Name:        "fnafify",
			Description: "ouuuuUUoOOuOUUHH SCARY?!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "the text to fnafify",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		"merge": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://fnaf.starchie.mom/?text=noob",
				},
			})
		},
		"fnafify": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			// TODO: refactor to helper function
			options := interaction.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			// end todo block
			text := optionMap["text"].StringValue()
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://fnaf.starchie.mom/?text=" + text,
				},
			})
		},
	}
)

func init() {
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
			handler(session, interaction)
		}
	})
}

func main() {
	session.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		log.Printf("logged in as: %v#%v", ready.User.Username, ready.User.Discriminator)
	})
	err := session.Open()
	if err != nil {
		log.Fatalf("cannot open the session: %v", err)
	}
	log.Println("adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", command)
		if err != nil {
			log.Panicf("cannot create '%v' command: %v", command.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// destruct discord session
	defer session.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("press ctrl+c to exit")
	<-stop

	if *RemoveCommands {
		log.Println("removing commands...")
		for _, v := range registeredCommands {
			err := session.ApplicationCommandDelete(session.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("gracefully shutting down.")
}
