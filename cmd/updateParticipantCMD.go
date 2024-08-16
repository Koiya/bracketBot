package cmd

import (
	"bracketBot/util"
	"github.com/bwmarrin/discordgo"
)

func UpdateParticipantCMD(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
	var name string
	var seed int64
	var discordUser string
	//var email string
	//var username string
	if opt, ok := optionMap["name"]; ok {
		name = opt.StringValue()
	}
	if opt, ok := optionMap["seed"]; ok {
		seed = opt.IntValue()
	}
	//if opt, ok := optionMap["email"]; ok {
	//	email = opt.StringValue()
	//}
	//if opt, ok := optionMap["username"]; ok {
	//	username = opt.StringValue()
	//}
	if opt, ok := optionMap["discord-user"]; ok {
		discordUser = opt.StringValue()
	}
	customOpt := util.Options{
		Name: name,
		Seed: int(seed),
		Misc: discordUser,
		//Email:    email,
		//Username: username,
	}
	message = util.UpdateParticipant(tourneyID, participantID, customOpt)
	cmd := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
	return s.InteractionRespond(i.Interaction, cmd)
}
