package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func ShowParticipantCMD(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	//Gets the params from the command
	options := i.ApplicationCommandData().Options[0].Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var participantID = optionMap["participant-id"].StringValue()
	var tourneyID = optionMap["tourney-id"].StringValue()
	data := util.FetchParticipant(tourneyID, participantID)
	name := data[0]
	ID := data[1]
	seed := data[2]
	misc := data[3]
	cmd := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("Participant Info"),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Name (Discord user)",
							Value:  fmt.Sprintf("%s (%s)", name, misc),
							Inline: true,
						},
						{
							Name:   "ID",
							Value:  ID,
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
	}
	return s.InteractionRespond(i.Interaction, cmd)
}
