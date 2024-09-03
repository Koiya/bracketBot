package cmd

import (
	"bracketBot/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
)

func UpdateTournamentState(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	util.SendRoleCheckMessage(s, i)

	//Gets the params from the command
	options := i.ApplicationCommandData().Options
	opt := options[0].Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opt))
	for _, o := range opt {
		optionMap[o.Name] = o
	}
	var tourneyID = optionMap["tourney-id"].StringValue()
	var states = optionMap["states"].StringValue()
	URL := fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/change_state.json", tourneyID)
	requestBody := fmt.Sprintf(`{
	"data": {
		"type" : "TournamentState",
		"attributes" : {
			"state":"%v"
			}
		}
	}`, states)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", URL, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", util.GetTOML("API.token"))

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Body:", string(body))

	if resp.StatusCode != 200 {
		var errorData util.ErrorWrapper
		err = json.Unmarshal(body, &errorData)
		//fmt.Println(resp.StatusCode)
		//fmt.Println(errorData)
		errDetail := errorData.Errors[0].Detail
		cmd := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", errDetail),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}
		return s.InteractionRespond(i.Interaction, cmd)
	}
	cmd := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Updated tournament state to %v", states),
		},
	}
	return s.InteractionRespond(i.Interaction, cmd)
}
