package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func ShowMatchCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//Gets the params from the command
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		var matchID = optionMap["match-id"].StringValue()
		var tourneyID = optionMap["tourney-id"].StringValue()
		data := util.FetchMatch(tourneyID, matchID)
		pOne := data[0]
		score := data[1]
		pTwo := data[2]
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: fmt.Sprintf("Matches in %v", tourneyID),
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Participant 1",
								Value:  pOne,
								Inline: true,
							},
							{
								Name:   "Score",
								Value:  score,
								Inline: true,
							},
							{
								Name:   "Participant 2",
								Value:  pTwo,
								Inline: true,
							},
						},
					},
				},
				AllowedMentions: &discordgo.MessageAllowedMentions{},
			},
		})
	}
	return cmd
}
