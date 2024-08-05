package util

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"log"
//	"net/http"
//	"strconv"
//)
//
///*
//	{
//	  "data": {
//	    "type": "Tournaments",
//	    "attributes": {
//	      "name": "Rainbow Six: Siege Tournament",
//	      "url": "r6siege_tournament_xx",
//	      "tournament_type": "single elimination",
//	      "game_name": "Tom Clancy%27s Rainbow Six Siege",
//	      "private": false,
//	      "starts_at": "string",
//	      "description": "string",
//	      "notifications": {
//	        "upon_matches_open": true,
//	        "upon_tournament_ends": true
//	      },
//	      "match_options": {
//	        "consolation_matches_target_rank": 3,
//	        "accept_attachments": false
//	      },
//	      "registration_options": {
//	        "open_signup": false,
//	        "signup_cap": 0,
//	        "check_in_duration": 0
//	      },
//	      "seeding_options": {
//	        "hide_seeds": false,
//	        "sequential_pairings": false
//	      },
//	      "station_options": {
//	        "auto_assign": false,
//	        "only_start_matches_with_assigned_stations": false
//	      },
//	      "group_stage_enabled": false,
//	      "group_stage_options": {
//	        "stage_type": "round robin",
//	        "group_size": 4,
//	        "participant_count_to_advance_per_group": 2,
//	        "rr_iterations": 1,
//	        "ranked_by": "",
//	        "rr_pts_for_match_win": 1,
//	        "rr_pts_for_match_tie": 0.5,
//	        "rr_pts_for_game_win": 0,
//	        "rr_pts_for_game_tie": 0,
//	        "split_participants": false
//	      },
//	      "double_elimination_options": {
//	        "split_participants": false,
//	        "grand_finals_modifier": ""
//	      },
//	      "round_robin_options": {
//	        "iterations": 2,
//	        "ranking": "",
//	        "pts_for_game_win": 1,
//	        "pts_for_game_tie": 0,
//	        "pts_for_match_win": 1,
//	        "pts_for_match_tie": 0.5
//	      },
//	      "swiss_options": {
//	        "rounds": 2,
//	        "pts_for_game_win": 1,
//	        "pts_for_game_tie": 0,
//	        "pts_for_match_win": 1,
//	        "pts_for_match_tie": 0.5
//	      },
//	      "free_for_all_options": {
//	        "max_participants": 4
//	      }
//	    }
//	  }
//	}
//*/
//func createTournament(tourneyID string) [4]string {
//	var results [4]string
//	//Generate random url link
//	charset := "abcdefghijklmnopqrstuvwxyz1234567890"
//
//	randomURL := ""
//	//Request to the API
//	tournament_type := []string{"single elimination", "double elimination", "round robin", "swiss", "free for all"}
//
//	requestBody := fmt.Sprintf(`
//		{
//		  "data": {
//			"type": "Tournaments",
//			"attributes": {
//			  "name": %v,
//			  "url": "",
//			  "tournament_type": %v,
//			  "game_name": %v ,
//			  "private": %b,
//			  "starts_at": "string",
//			  "description": "",
//			  "notifications": {
//				"upon_matches_open": true,
//				"upon_tournament_ends": true
//			  },
//			  "match_options": {
//				"consolation_matches_target_rank": 3,
//				"accept_attachments": false
//			  },
//			  "registration_options": {
//				"open_signup": false,
//				"signup_cap": 0,
//				"check_in_duration": %d
//			  },
//			  "seeding_options": {
//				"hide_seeds": false,
//				"sequential_pairings": false
//			  },
//			  "station_options": {
//				"auto_assign": false,
//				"only_start_matches_with_assigned_stations": false
//			  },
//			  "group_stage_enabled": false,
//			  "group_stage_options": {
//				"stage_type": "round robin",
//				"group_size": 4,
//				"participant_count_to_advance_per_group": 2,
//				"rr_iterations": 1,
//				"ranked_by": "",
//				"rr_pts_for_match_win": 1,
//				"rr_pts_for_match_tie": 0.5,
//				"rr_pts_for_game_win": 0,
//				"rr_pts_for_game_tie": 0,
//				"split_participants": false
//			  },
//			  "double_elimination_options": {
//				"split_participants": false,
//				"grand_finals_modifier": ""
//			  },
//			  "round_robin_options": {
//				"iterations": 2,
//				"ranking": "",
//				"pts_for_game_win": 1,
//				"pts_for_game_tie": 0,
//				"pts_for_match_win": 1,
//				"pts_for_match_tie": 0.5
//			  },
//			  "swiss_options": {
//				"rounds": 2,
//				"pts_for_game_win": 1,
//				"pts_for_game_tie": 0,
//				"pts_for_match_win": 1,
//				"pts_for_match_tie": 0.5
//			  },
//			  "free_for_all_options": {
//				"max_participants": 4
//			  }
//			}
//		  }
//		}
//		`)
//
//	var URL string
//	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments.json")
//	client := &http.Client{}
//	req, err := http.NewRequest("POST", URL, bytes.NewBuffer([]byte(requestBody)))
//	if err != nil {
//		log.Fatal(err)
//	}
//	req.Header.Set("Accept", "application/json")
//	req.Header.Set("Content-Type", "application/vnd.api+json")
//	req.Header.Set("Authorization-Type", "v1")
//	req.Header.Set("Authorization", GetTOML("API.token"))
//
//	resp, err := client.Do(req)
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	var data AllParticipantWrapper
//	err = json.Unmarshal(body, &data)
//	if err != nil {
//		log.Fatal("Error", err.Error())
//	}
//
//	if resp.StatusCode != 200 {
//		fmt.Println(resp.StatusCode)
//		results[0] = "Error getting data"
//		return results
//	}
//	for _, v := range data.Data {
//		results[0] += fmt.Sprintf("%v (%v)\n\n\n", v.Attributes.Name, v.Id)
//		results[1] += v.Attributes.Misc + "\n\n\n"
//		results[2] += strconv.Itoa(v.Attributes.Seed) + "\n\n\n"
//	}
//	//fmt.Println(data.Data)
//
//	return results
//}
