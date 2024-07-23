package cmd

import (
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate,) {
	_, err := s.ApplicationCommandBulkOverwrite(*appID, *guildID, []*discordgo.ApplicationCommand{
		{
			Name:        "hello-world",
			Description: "Showcase of a basic slash command",
		},
	})
	if err != nil {
		// Handle the error
	}
	data := i.ApplicationCommandData()
	switch data.Name {
	case "hello-world":
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: "Hello world!",
				},
			},
		)
		if err != nil {
			// Handle the error
		}
	}
})
}
