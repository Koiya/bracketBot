package cmd

import (
	"bracketBot/util"
	"github.com/bwmarrin/discordgo"
)

func UpdateMatchCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		var scores = optionMap["scores"].StringValue()
		var matchID = optionMap["match-id"].StringValue()
		var advancing = optionMap["advancing"].BoolValue()

		customOpt := util.MatchOptions{
			ParticipantID: participantID,
			Score:         scores,
			Advancing:     advancing,
		}
		if !util.RoleCheck(i) {
			message = "You don't have permission to use this command"
			goto Skip
		}
		message = util.UpdateMatch(tourneyID, matchID, customOpt)
	Skip:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
	return cmd
}
