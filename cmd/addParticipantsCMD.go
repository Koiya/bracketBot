package cmd

import (
	"bracketBot/util"
	"github.com/bwmarrin/discordgo"
)

func AddParticipantsCMD(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	util.SendRoleCheckMessage(s, i)
	var message string

	//Gets the params from the command
	options := i.ApplicationCommandData().Options[0].Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var tourneyID = optionMap["tourney-id"].StringValue()
	var name = optionMap["name"].StringValue()
	//Discord int is 64bit
	var seed int64
	var discordUser = optionMap["discord-user"].StringValue()
	// Uncomment if using email and challonge username
	//var email string
	//var username string

	// Get the value from the option map.
	// When the option exists, ok = true
	if opt, ok := optionMap["seed"]; ok {
		seed = opt.IntValue()
	}

	//if opt, ok := optionMap["email"]; ok {
	//	email = opt.StringValue()
	//}
	//if opt, ok := optionMap["username"]; ok {
	//	username = opt.StringValue()
	//}
	//seed must be converted to 32bit after getting value
	customOpt := util.Options{
		Name: name,
		Seed: int(seed),
		Misc: discordUser,
		//Email:    email,
		//Username: username,
	}
	message = util.AddParticipants(tourneyID, customOpt)
	cmd := &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
	return s.InteractionRespond(i.Interaction, cmd)
}
