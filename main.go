package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"os/signal"
)

var (
	s       *discordgo.Session
	GuildID string
)

func init() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		fmt.Println("Error ", err.Error())
		return
	}
	botToken := config.Get("Bot.token").(string)
	GuildID = config.Get("Bot.guild_id").(string)
	s, err = discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session,", err.Error())
		return
	}

}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "ping",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Ping pong.",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong",
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v  GuildID: %v", s.State.User.Username, s.State.User.Discriminator, GuildID)
	})

	// add a event handler
	if err := s.Open(); err != nil {
		fmt.Println("Error opening connection,", err.Error())
		return
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	// keep cmd running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
