package cmd

import (
	"bracketBot/util"
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"strings"
)

func UpdateTournament(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
	var name = optionMap["name"].StringValue()
	var gameName = optionMap["game_name"].StringValue()
	var tournamentType = strings.ToLower(optionMap["tournament_type"].StringValue())
	var startTime = optionMap["start_time"].StringValue()

	//Request to the API

	requestBody := fmt.Sprintf(`
		{
		  "data": {
			"type": "Tournaments",
			"attributes": {
			  "name": "%v",
			  "tournament_type": "%v",
			  "game_name": "%v" ,
			  "private": "true",
			  "starts_at": "%v",
			  "description": "",
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
				"check_in_duration": 15
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
		`, name, tournamentType, gameName, startTime)
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v.json", tourneyID)
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer([]byte(requestBody)))
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

	if resp.StatusCode != 201 {
		fmt.Println(resp.StatusCode)
		cmd := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error occurred. Please try again.",
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
