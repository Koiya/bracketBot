package cmd

import (
	"github.com/bwmarrin/discordgo"
)

func Ping() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong",
			},
		})
	}
	return cmd
}
