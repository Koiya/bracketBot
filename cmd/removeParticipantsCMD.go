package cmd

import (
	"bracketBot/util"
	"github.com/bwmarrin/discordgo"
)

func RemoveParticipantCMD(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	util.SendRoleCheckMessage(s, i)
	var message string

	//Gets the params from the command
	options := i.ApplicationCommandData().Options[0].Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var tourneyID = optionMap["tourney-id"].StringValue()
	var participantID = optionMap["participant-id"].StringValue()
	message = util.RemoveParticipants(tourneyID, participantID)
	cmd := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
	return s.InteractionRespond(i.Interaction, cmd)
}
