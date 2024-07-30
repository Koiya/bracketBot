package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func ShowTournamentCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//Gets the params from the command
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		var tourneyID = optionMap["tourney-id"].StringValue()
		data := util.FetchATournament(tourneyID)
		name := data[0]
		game := data[1]
		tourneyType := data[2]
		url := data[3]
		imageURL := data[4]

		//NEED TO CONVERT SVG TO PNG OR JPG SOMEHOW

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: fmt.Sprintf("%s  -  %s  - %s", name, tourneyType, game),
						URL:   url,
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Link to bracket",
								Value:  imageURL,
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
