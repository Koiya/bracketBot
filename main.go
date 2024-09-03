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
		//CREATE COMMANDS
		{
			Name:        "create",
			Description: "Create tourney/participant",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "tournament",
					Description: "Create a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Single Elimination",
									Value: "single elimination",
								},
								{
									Name:  "Double Elimination",
									Value: "double elimination",
								},
								{
									Name:  "Round Robin",
									Value: "round robin",
								},
								{
									Name:  "Swiss",
									Value: "swiss",
								},
								{
									Name:  "Free for all",
									Value: "free for all",
								},
							},
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
				},
				{
					Name:        "participant",
					Description: "Add a participant to a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "discord-user",
							Description: "Discord username of the participant",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "seed",
							Description: "Seeding of the participant",
							Required:    false,
						},
					},
				},
			},
		},
		//REMOVE COMMANDS
		{
			Name:        "remove",
			Description: "Remove a tourney/participant",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "tournament",
					Description: "Remove a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "tourney-id",
							Description: "Input the tournament id that you want to delete",
							Required:    true,
						},
					},
				},
				{
					Name:        "participant",
					Description: "Removes a participant from a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
			},
		},
		//UPDATE COMMANDS
		{
			Name:        "update",
			Description: "Update tourney/participant/match/state",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "tournament",
					Description: "Update tournament information",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "tourney-id",
							Description: "Input the tournament id that you want to delete",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Name of the tournament",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "tournament_type",
							Description: "Single Elimination, Double Elimination, Round Robin, Swiss, Free for all",
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Single Elimination",
									Value: "single elimination",
								},
								{
									Name:  "Double Elimination",
									Value: "double elimination",
								},
								{
									Name:  "Round Robin",
									Value: "round robin",
								},
								{
									Name:  "Swiss",
									Value: "swiss",
								},
								{
									Name:  "Free for all",
									Value: "free for all",
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "game_name",
							Description: "Name of the game",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "start_time",
							Description: "Input date and time",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "check_in",
							Description: "Check in window",
							Required:    false,
						},
					},
				},
				{

					Name:        "tournamentstate",
					Description: "Update tournament's state",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
				},
				{
					Name:        "participant",
					Description: "Update a participant from a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
							Name:        "discord_user",
							Description: "Discord username of the participant",
							Required:    false,
						},
					},
				},
				{
					Name:        "match",
					Description: "Update a match from a tournament",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
			},
		},
		//SHOW COMMANDS
		{
			Name:        "show",
			Description: "Display matches/tournaments/participants",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "all",
					Description: "Display all matches/tournaments/participants",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "tournament",
							Description: "Display all tournaments ",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
						{
							Name:        "match",
							Description: "Display all matches in a tournament",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Type:        discordgo.ApplicationCommandOptionString,
									Name:        "tourney-id",
									Description: "Input ID of the tournament",
									Required:    true,
								},
								{
									Type:        discordgo.ApplicationCommandOptionString,
									Name:        "states",
									Description: "Check all open or pending state of the matches",
									Required:    true,
									Choices: []*discordgo.ApplicationCommandOptionChoice{
										{
											Name:  "Open",
											Value: "open",
										},
										{
											Name:  "Pending",
											Value: "pending",
										},
										{
											Name:  "Completed",
											Value: "complete",
										},
									},
								},
							},
						},
						{
							Name:        "participant",
							Description: "Display all participants in a tournament ",
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
					Name:        "tournament",
					Description: "Display a tournament with given ID",
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
				{
					Name:        "match",
					Description: "Display a match with given IDs",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
					Name:        "participant",
					Description: "Display a participant with given IDs",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
			},
		},
		//UTILITY COMMANDS
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
		//CREATE HANDLER
		"create": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
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
				if err := cmd.AddParticipantsCMD(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			}
		},
		//REMOVE HANDLER
		"remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
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
				if err := cmd.RemoveParticipantCMD(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			}

		},
		//UPDATE HANDLER
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
				if err := cmd.UpdateParticipantCMD(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			case "match":
				if err := cmd.UpdateMatchCMD(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			}
		},
		"show": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			switch options[0].Name {
			case "tournament":
				if err := cmd.ShowTournamentCMD(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			case "match":
				if err := cmd.ShowMatchCMD(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			case "participant":
				if err := cmd.ShowParticipantCMD(s, i); err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error occurred when using command. Please try again later.",
						},
					})
					fmt.Println(err)
				}
			case "all":
				options = options[0].Options
				switch options[0].Name {
				case "tournament":
					if err := cmd.ShowAllTournamentsCMD(s, i); err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Error occurred when using command. Please try again later.",
							},
						})
						fmt.Println(err)
					}
				case "match":
					if err := cmd.ShowAllMatchesCMD(s, i); err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Error occurred when using command. Please try again later.",
							},
						})
						fmt.Println(err)
					}
				case "participant":
					if err := cmd.ShowAllParticipantsCMD(s, i); err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Error occurred when using command. Please try again later.",
							},
						})
						fmt.Println(err)
					}
				}
			}
		},
		"rollcall": cmd.RollCallCMD(),
	}
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"rc_join":  cmd.RCJoinComponent(),
		"rc_close": cmd.RCCloseComponent(),
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
