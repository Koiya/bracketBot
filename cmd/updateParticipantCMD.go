package cmd

import (
	"bracketBot/util"
	"github.com/bwmarrin/discordgo"
)

func UpdateParticipantCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		var name string
		var seed int64
		var misc string
		var email string
		var username string
		if opt, ok := optionMap["name"]; ok {
			name = opt.StringValue()
		}
		if opt, ok := optionMap["seed"]; ok {
			seed = opt.IntValue()
		}
		if opt, ok := optionMap["email"]; ok {
			email = opt.StringValue()
		}
		if opt, ok := optionMap["username"]; ok {
			username = opt.StringValue()
		}
		if opt, ok := optionMap["misc"]; ok {
			misc = opt.StringValue()
		}
		customOpt := util.Options{
			Name:     name,
			Seed:     int(seed),
			Misc:     misc,
			Email:    email,
			Username: username,
		}
		if !util.RoleCheck(i) {
			message = "You don't have permission to use this command"
			goto Skip
		}
		message = util.UpdateParticipant(tourneyID, participantID, customOpt)
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