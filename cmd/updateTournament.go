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
		if checkIn < 14 {
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
			  "private": "true",
			  "description": "",
			  "start_at": "%v",
			  "notifications": {
				"upon_matches_open": true,
				"upon_tournament_ends": true
			  },
			  "match_options": {
				"consolation_matches_target_rank": 3,
				"accept_attachments": false
			  },
			  "registration_options": {
				"open_signup": false,
				"signup_cap": 5,
				"check_in_duration" : %d
			  },
			  "seeding_options": {
				"hide_seeds": false,
				"sequential_pairings": false
			  },
			  "station_options": {
				"auto_assign": false,
				"only_start_matches_with_assigned_stations": false
			  },
			  "group_stage_enabled": false,
			  "group_stage_options": {
				"stage_type": "round robin",
				"group_size": 4,
				"participant_count_to_advance_per_group": 2,
				"rr_iterations": 1,
				"ranked_by": "match wins",
				"rr_pts_for_match_win": 1,
				"rr_pts_for_match_tie": 0.5,
				"rr_pts_for_game_win": 0,
				"rr_pts_for_game_tie": 0,
				"split_participants": false
			  },
			  "double_elimination_options": {
				"split_participants": false,
				"grand_finals_modifier": ""
			  },
			  "round_robin_options": {
				"iterations": 2,
				"ranking": "match wins",
				"pts_for_game_win": 1,
				"pts_for_game_tie": 0,
				"pts_for_match_win": 1,
				"pts_for_match_tie": 0.5
			  },
			  "swiss_options": {
				"rounds": 2,
				"pts_for_game_win": 1,
				"pts_for_game_tie": 0,
				"pts_for_match_win": 1,
				"pts_for_match_tie": 0.5
			  },
			  "free_for_all_options": {
				"max_participants": 4
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
		log.Fatal(err)
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
		log.Fatal(err)
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
