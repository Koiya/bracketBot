package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func errorResponse(err error) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Error has occured", err.Error())
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "error",
			},
		})
	}
	return cmd
}

func ShowTournaments() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
