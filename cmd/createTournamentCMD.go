package cmd

import (
	"bracketBot/util"
	"github.com/bwmarrin/discordgo"
)

func CreateTournamentCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var message string

		//Gets the params from the command
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		//var tourneyID = optionMap["tourney-id"].StringValue()
		//var name = optionMap["name"].StringValue()
		//var seed int64
		//var misc string
		//var email string
		//var username string

		// Get the value from the option map.
		// When the option exists, ok = true
		//if opt, ok := optionMap["seed"]; ok {
		//	seed = opt.IntValue()
		//}
		//if opt, ok := optionMap["email"]; ok {
		//	email = opt.StringValue()
		//}
		//if opt, ok := optionMap["username"]; ok {
		//	username = opt.StringValue()
		//}
		//if opt, ok := optionMap["misc"]; ok {
		//	misc = opt.StringValue()
		//}
		//customOpt := util.Options{
		//	Seed:     int(seed),
		//	Misc:     misc,
		//	Email:    email,
		//	Username: username,
		//}
		if !util.RoleCheck(i) {
			message = "You don't have permission to use this command"
			goto Skip
		}
	Skip:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
	/*	url := fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%d/participants.json", TourneyID)
		d := Participant{
			name,
			seed,
			misc,
			email,
			username,
		}
		b, err := json.Marshal(d)*/

	return cmd
}
