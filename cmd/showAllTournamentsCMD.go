package cmd

import (
	"bracketBot/util"
	"github.com/bwmarrin/discordgo"
)

func ShowAllTournamentsCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := util.FetchAllTournaments()
		name := data[0]
		Game := data[1]
		Type := data[2]
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Tournaments",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Name",
								Value:  name,
								Inline: true,
							},
							{
								Name:   "Game",
								Value:  Game,
								Inline: true,
							},
							{
								Name:   "Type",
								Value:  Type,
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
