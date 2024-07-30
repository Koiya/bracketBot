package cmd

import (
	"bracketBot/util"
	"github.com/bwmarrin/discordgo"
)

func RemoveParticipantCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var message string

		//Gets the params from the command
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		var tourneyID = optionMap["tourney-id"].StringValue()
		var participantID = optionMap["participant-id"].StringValue()
		if !util.RoleCheck(i) {
			message = "You don't have permission to use this command"
			goto Skip
		}
		message = util.RemoveParticipants(tourneyID, participantID)
	Skip:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
	return cmd
}
