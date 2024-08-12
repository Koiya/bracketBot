package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func RollCallCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//Gets the params from the command
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		var tourneyID = optionMap["tourney-id"].StringValue()
		data := util.FetchATournament(tourneyID)
		tourneyName := data[0]
		fmt.Println(i.Member.User.Username)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Rollcall for " + tourneyName,
						Description: "has started rollcall",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Start time",
								Value:  "TIME GOES HERE",
								Inline: true,
							},
						},
					},
				},
				AllowedMentions: &discordgo.MessageAllowedMentions{},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Emoji: &discordgo.ComponentEmoji{
									Name: "➡️",
								},
								Label:    "Join",
								Style:    discordgo.SuccessButton,
								CustomID: "rc_join",
							},
							discordgo.Button{
								Emoji: &discordgo.ComponentEmoji{
									Name: "✖️",
								},
								Label:    "Close",
								Style:    discordgo.DangerButton,
								CustomID: "rc_leave",
							},
							discordgo.Button{
								Emoji: &discordgo.ComponentEmoji{
									Name: "🔗",
								},
								Label: "URL",
								Style: discordgo.LinkButton,
								URL:   "https://challonge.com/" + tourneyID,
							},
						},
					},
				},
			},
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	return cmd
}

func RCJoinComponent() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmt.Println(i.Message.ChannelID)
		fmt.Println(i.Message.ID)
		//optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		//for _, opt := range options {
		//	optionMap[opt.Name] = opt
		//}
		//var tourneyID = optionMap["tourney-id"].StringValue()
		name := i.Member.User.GlobalName
		discordUser := i.Member.User.Username
		fmt.Println(name + " " + discordUser)
		//customOpt := util.Options{
		//	Name: name,
		//	Seed: 1,
		//	Misc: misc,
		//}
		//message := util.AddParticipants(tourneyID, customOpt)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "TEST JOIN",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return cmd
}