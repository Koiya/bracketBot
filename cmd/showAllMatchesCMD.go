package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func ShowAllMatchesCMD(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	//Gets the params from the command
	options := i.ApplicationCommandData().Options[0].Options[0].Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var tourneyID = optionMap["tourney-id"].StringValue()
	tourneyData := util.FetchATournament(tourneyID)
	tourneyName := tourneyData[0]
	data := util.FetchAllMatches(tourneyID)
	pOne := data[0]
	score := data[1]
	pTwo := data[2]
	cmd := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("Matches in %v", tourneyName),
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
	}
	return s.InteractionRespond(i.Interaction, cmd)
}
