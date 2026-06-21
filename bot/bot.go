package main

import (
	"flag"
	"fmt"
	"log"
	nurl "net/url"
	"os"
	"os/signal"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var (
	fnafAPI        = "localhost:9638"
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
)

var session *discordgo.Session

func init() { flag.Parse() }

func init() {
	// load env
	token := os.Getenv("DISCORD_TOKEN")

	// setup session bot session
	var err error
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "the image or @user to fnafify with",
					MinLength:   new(1),
					MaxLength:   500,
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "attachment",
					Description: "the image to fnafify with",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "preset",
					Description: "the image to fnafify with",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "not implemented yet",
							Value: "not_implemented",
						},
					},
					Required: false,
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
			data := interaction.ApplicationCommandData()
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(data.Options))
			for _, opt := range data.Options {
				optionMap[opt.Name] = opt
			}
			// end todo block
			text := ""
			if optionMap["text"] != nil {
				text = optionMap["text"].StringValue()
			}

			url := ""
			if optionMap["url"] != nil {
				url = optionMap["url"].StringValue()
			}

			preset := ""
			if optionMap["preset"] != nil {
				url = optionMap["preset"].StringValue()
			}

			attachment := optionMap["attachment"]

			log.Printf("text: %s, url: %s, preset: %s, attachment: %v", text, url, preset, attachment)

			content := fmt.Sprintf("https://fnaf.starchie.mom/?text=%s", nurl.QueryEscape(text))
			if url != "" {
				matched, err := regexp.MatchString(`^<@!?(\d+)>$`, url)
				if err == nil && matched {
					// sanitize stupid string "<@123456789>" -> "123456789")
					re := regexp.MustCompile(`\d+`)
					userID := re.FindString(url)

					// get the user object from the session cache or api
					user, err := session.User(userID)
					if err != nil {
						log.Printf("error: could not get user with id %s: %v", userID, err)
					} else {
						// get the user's avatar url only if we successfully got the user
						avatarURL := user.AvatarURL("512")
						url = avatarURL
					}
				}
				content += fmt.Sprintf("&url=%s", nurl.QueryEscape(url))
			} else if attachment != nil {
				var attachmentURL string
				for _, attachment := range data.Resolved.Attachments {
					attachmentURL = nurl.QueryEscape(attachment.URL)
					break // only get first attachment for now
				}
				content += fmt.Sprintf("&url=%s", attachmentURL)
			}
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
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
