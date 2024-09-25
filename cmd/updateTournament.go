package cmd

import (
	"bracketBot/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func UpdateTournament(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	//Check if user has permissions if not it'll return a different message.
	util.SendRoleCheckMessage(s, i)

	//Gets the params from the command
	options := i.ApplicationCommandData().Options
	opt := options[0].Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opt))
	for _, o := range opt {
		optionMap[o.Name] = o
	}
	var tourneyID = optionMap["tourney-id"].StringValue()
	var tourneyInfo = util.FetchATournament(tourneyID)
	var name string
	var gameName string
	var tournamentType string
	var startTime string
	var checkIn int
	fmt.Println(tourneyInfo)

	if opt, ok := optionMap["name"]; ok {
		name = opt.StringValue()
	} else {
		name = tourneyInfo[0]
	}

	if opt, ok := optionMap["game_name"]; ok {
		gameName = opt.StringValue()
	} else {
		gameName = tourneyInfo[1]
	}

	if opt, ok := optionMap["tournament_type"]; ok {
		tournamentType = strings.ToLower(opt.StringValue())
	} else {
		tournamentType = tourneyInfo[2]
	}

	if opt, ok := optionMap["start_time"]; ok {
		startTime = opt.StringValue()
	} else {
		startTime = tourneyInfo[5]
	}
	fmt.Println(name, gameName, tournamentType, startTime)
	if opt, ok := optionMap["check_in"]; ok {
		checkIn = int(opt.IntValue())
		if checkIn <= 14 {
			cmd := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Check in time must be at least 15 minutes",
				},
			}
			return s.InteractionRespond(i.Interaction, cmd)
		}
	} else {
		checkIn, _ = strconv.Atoi(tourneyInfo[6])
		if checkIn <= 14 {
			checkIn = 15
		}
	}

	//Request to the API

	requestBody := fmt.Sprintf(`
		{
		  "data": {
			"type": "Tournaments",
			"attributes": {
			  "name" : "%v",
			  "game_name" : "%v",
			  "url" : "%v",
              "tournament_type" : "%v",
			  "start_at": "%v",
			  "registration_options": {
				"check_in_duration" : %d
			  }
			}
		  }
		}
		`, name, gameName, tourneyID, tournamentType, startTime, checkIn)
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v.json", tourneyID)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", URL, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", util.GetTOML("API.token"))

	resp, err := client.Do(req)
	defer resp.Body.Close()

	//replace _ to body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("response Body:", string(body))

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		var errorData util.ErrorWrapperArray
		err = json.Unmarshal(body, &errorData)
		//fmt.Println(resp.StatusCode
		errDetail := errorData.Errors[0].Detail[0]
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
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       fmt.Sprintf("%v has been updated", name),
					Description: gameName,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Tournament Type",
							Value:  tournamentType,
							Inline: true,
						},
						{
							Name:   "Start time",
							Value:  startTime,
							Inline: true,
						},
					},
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "🔗",
							},
							Label: "URL",
							Style: discordgo.LinkButton,
							URL:   "https://challonge.com/" + tourneyID,
						},
					},
				},
			},
		},
	}
	return s.InteractionRespond(i.Interaction, cmd)
}
