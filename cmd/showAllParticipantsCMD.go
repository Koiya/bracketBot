package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func ShowAllParticipantsCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//Gets the params from the command
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		var tourneyID = optionMap["tourney-id"].StringValue()
		data := util.FetchAllParticipants(tourneyID)
		name := data[0]
		misc := data[1]
		seed := data[2]
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: fmt.Sprintf("Participants in %v", tourneyID),
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Name",
								Value:  name,
								Inline: true,
							},
							{
								Name:   "Misc",
								Value:  misc,
								Inline: true,
							},
							{
								Name:   "Seed",
								Value:  seed,
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
