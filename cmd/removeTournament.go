package cmd

import (
	"bracketBot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
)

func RemoveTournament(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	//Check if user has permissions if not it'll return a different message.
	if !util.RoleCheck(i) {
		cmd := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have permission to use this command!",
			},
		}
		return s.InteractionRespond(i.Interaction, cmd)
	}
	//Gets the params from the command
	options := i.ApplicationCommandData().Options
	opt := options[0].Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opt))
	for _, o := range opt {
		optionMap[o.Name] = o
	}
	var tourneyID = optionMap["tourney-id"].StringValue()

	//Request to the API
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v.json", tourneyID)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", util.GetTOML("API.token"))

	resp, err := client.Do(req)
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 204 {
		fmt.Println(resp.StatusCode)
		cmd := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error deleting a tournament",
			},
		}
		return s.InteractionRespond(i.Interaction, cmd)
	}

	cmd := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Tournament has been deleted %v", tourneyID),
		},
	}
	return s.InteractionRespond(i.Interaction, cmd)
}
