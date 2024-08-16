package main

import (
	"bracketBot/cmd"
	"bracketBot/util"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

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

		//revise all the commands to use this format ex: update participant/tournament/match, show all (participant/tournament/match)
		{
			Name:        "subcommands",
			Description: "Subcommands and command groups example",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "subcommand-group",
					Description: "Subcommands group",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "nested-subcommand",
							Description: "Nested subcommand",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Type:        discordgo.ApplicationCommandOptionString,
									Name:        "tourney-id",
									Description: "Input ID of the tournament",
									Required:    true,
								},
							},
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				},
				{
					Name:        "subcommand",
					Description: "Top-level subcommand",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "create",
			Description: "Create tourney/participant",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "tournament",
					Description: "Create a tournament",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Name of the tournament",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "tournament_type",
							Description: "Single Elimination, Double Elimination, Round Robin, Swiss, Free for all",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "game_name",
							Description: "Name of the game",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "start_time",
							Description: "Input date and time",
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "participant",
					Description: "Insert a participant in a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "remove",
			Description: "Remove a tourney/participant",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "tournament",
					Description: "Remove a tournament",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "tourney-id",
							Description: "Input the tournament id that you want to delete",
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "participant",
					Description: "Insert a participant in a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "update",
			Description: "update a tourney/participant",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "tournamentstate",
					Description: "update a tournament",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "tourney-id",
							Description: "Input the tournament id that you want to delete",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "states",
							Description: "State of the tournament",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Process Checkin",
									Value: "process_checkin",
								},
								{
									Name:  "Start Group Stage",
									Value: "start_group_stage",
								},
								{
									Name:  "Finalize Group Stage",
									Value: "finalize_group_stage",
								},
								{
									Name:  "Reset Group Stage",
									Value: "reset_group_stage",
								},
								{
									Name:  "Start",
									Value: "start",
								},
								{
									Name:  "Finalize",
									Value: "finalize",
								},
								{
									Name:  "Reset",
									Value: "reset",
								},
								{
									Name:  "Open Predictions",
									Value: "open_predictions",
								},
							}},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "participant",
					Description: "Insert a participant in a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		//Tournaments
		{
			Name:        "showalltournaments",
			Description: "show all tournament",
		},
		{
			Name:        "showtournament",
			Description: "Show tournament with given ID",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
			},
		},
		//Match
		{
			Name:        "showmatch",
			Description: "Show a match from a tournament",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "match-id",
					Description: "Input ID of the match",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
			},
		},
		{
			Name:        "showmatches",
			Description: "Show all matches from a tournament",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
			},
		},
		{
			Name:        "updatematch",
			Description: "Update a match from a tournament",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "match-id",
					Description: "Input ID of the match",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "participant-id",
					Description: "Input ID of the participant",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "scores",
					Description: "Input scores Ex: 3 or 3,0,3",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "advancing",
					Description: "Advancing if won all set",
					Required:    true,
				},
			},
		},
		//Participants
		{
			Name:        "showparticipants",
			Description: "Show participants from a tournament",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
			},
		},
		{
			Name:        "showparticipant",
			Description: "Show participants from a tournament",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "participant-id",
					Description: "Input ID of the participant",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
			},
		},
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
					Description: "Discord username of the participant",
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
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "participant-id",
					Description: "Must be the ID of the participant",
					Required:    true,
				},
			},
		},
		{
			Name:        "updateparticipant",
			Description: "Update a participant from a tournament",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "participant-id",
					Description: "Must be the ID of the participant",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Name of the participant",
					Required:    false,
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
					Description: "Discord username of the participant",
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
			Name:        "rollcall",
			Description: "Start up match request and let player join off a button",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tourney-id",
					Description: "Input ID of the tournament",
					Required:    true,
				},
			},
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
		"subcommands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			content := ""

			// As you can see, names of subcommands (nested, top-level)
			// and subcommand groups are provided through the arguments.
			switch options[0].Name {
			case "subcommand":
				content = "The top-level subcommand is executed. Now try to execute the nested one."
			case "subcommand-group":
				options = options[0].Options
				fmt.Println(options)
				switch options[0].Name {
				case "nested-subcommand":
					content = "Nice, now you know how to execute nested commands too"
				default:
					content = "Oops, something went wrong.\n" +
						"Hol' up, you aren't supposed to see this message."
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"create": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			// As you can see, names of subcommands (nested, top-level)
			// and subcommand groups are provided through the arguments.
			switch options[0].Name {
			case "tournament":
				if err := cmd.CreateTournamentCMD(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			case "participant":
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "TEST",
					},
				})
			}

		},
		"remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			// As you can see, names of subcommands (nested, top-level)
			// and subcommand groups are provided through the arguments.
			switch options[0].Name {
			case "tournament":
				if err := cmd.RemoveTournament(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			case "participant":
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "TEST",
					},
				})
			}

		},
		"update": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			switch options[0].Name {
			case "tournament":
				if err := cmd.UpdateTournament(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			case "tournamentstate":
				if err := cmd.UpdateTournamentState(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			case "participant":
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "TEST",
					},
				})
			}
		},
		"showalltournaments": cmd.ShowAllTournamentsCMD(),
		"showtournament":     cmd.ShowTournamentCMD(),
		"showparticipants":   cmd.ShowAllParticipantsCMD(),
		"showparticipant":    cmd.ShowParticipantCMD(),
		"addparticipant":     cmd.AddParticipantsCMD(),
		"removeparticipant":  cmd.RemoveParticipantCMD(),
		"updateparticipant":  cmd.UpdateParticipantCMD(),
		"showmatches":        cmd.ShowAllMatchesCMD(),
		"showmatch":          cmd.ShowMatchCMD(),
		"updatematch":        cmd.UpdateMatchCMD(),
		"rollcall":           cmd.RollCallCMD(),
	}
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"rc_join": cmd.RCJoinComponent(),
	}
)

func init() {
}

func main() {
	var RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v  GuildID: %v", s.State.User.Username, s.State.User.Discriminator, GuildID)
	})
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:

			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
	// add a event handler
	if err := s.Open(); err != nil {
		fmt.Println("Error opening connection,", err.Error())
		return
	}
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	fmt.Println("Adding commands...")
	for i, v := range commands {
		fmt.Println("Added " + v.Name)
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
	if *RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
	log.Println("Gracefully shutting down.")
}
