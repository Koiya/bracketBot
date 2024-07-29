package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Participant struct {
	Name     string
	Seed     int
	Misc     string
	Email    string
	Username string
}

func init() {

}

func AddParticipantsCMD() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		//Gets the params from the command
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		//margs := make([]interface{}, 0, len(options))
		//msgformat := "You learned how to use command options! " +
		//	"Take a look at the value(s) you entered:\n"

		var tourneyID = optionMap["tourney-id"].StringValue()
		var name = optionMap["name"].StringValue()
		util.AddParticipants(tourneyID, name)

		// Get the value from the option map.
		// When the option exists, ok = true
		//if opt, ok := optionMap["seed"]; ok {
		//	var seed = opt.IntValue()
		//}
		//if opt, ok := optionMap["email"]; ok {
		//	var email = opt.StringValue()
		//}
		//if opt, ok := optionMap["username"]; ok {
		//	var username = opt.StringValue()
		//}
		//if opt, ok := optionMap["misc"]; ok {
		//	var misc = opt.StringValue()
		//}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"Added %v to %v", name, tourneyID,
				),
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

//https://api.challonge.com/v2.1/tournaments/{tourneyid}/participants.json
