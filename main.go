package main

import (
	"bracketBot/cmd"
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

//TODO
// Get one tournament
// Create Tournament
// Update a Tournament
// Delete a tournament(admin)
//
// show bracket
//participant
// -add
// -remove

var (
	s       *discordgo.Session
	GuildID string
	err     error
)

func init() {
	botToken := util.GetTOML("Bot.token")
	GuildID = util.GetTOML("Bot.guild_id")
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

		//Tournaments
		{
			Name:        "showalltournaments",
			Description: "show all tournament",
		},
		{
			Name:        "updatetournament",
			Description: "Updates a tournament with options passed",
		},

		//Participants
		{
			Name:        "addparticipant",
			Description: "Add a participant to a tournament",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Name of the participant",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "seed",
					Description: "Seeding of the participant",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "misc",
					Description: "Description of the participant",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "email",
					Description: "Email of the participant",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "Challonge username of the participant",
					Required:    false,
				},
			},
		},
		{
			Name:        "removeparticipant",
			Description: "Removes a participant from a tournament",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var message string
			if len(i.Member.Roles) == 0 {
				message = "No roles!"
			}
			for _, value := range i.Member.Roles {
				message += value + " "
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message,
				},
			})
		},
		"showalltournaments": cmd.ShowAllTournamentsCMD(),
		"addparticipant":     cmd.AddParticipantsCMD(),
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
