package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type MatchOptions struct {
	ParticipantID string
	Score         string
	Advancing     bool
}

func UpdateMatch(tourneyID, matchID string, opt MatchOptions) string {

	//Request to the API
	//Check if any other characters other than numbers or comma then return if true
	match, _ := regexp.MatchString("^[0-9,]*$", opt.Score)
	if !match {
		return "Score format incorrect. Numbers and commas only."
	}
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/matches/%v.json", tourneyID, matchID)
	requestBody := fmt.Sprintf(`{
	"data": {
		"type" : "Match",
		"attributes" : {
			"match": [
				{
					"participant_id": "%v",
					"score_set": "%v",
					"rank": 1,
					"advancing": %t
				}
			],
			"tie": false
			}
		}
	}`, opt.ParticipantID, opt.Score, opt.Advancing)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", URL, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", GetTOML("API.token"))

	resp, err := client.Do(req)
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(body))

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return "Error updating match. Please check parameters."
	}
	tourneyData := FetchATournament(tourneyID)
	participantData := FetchParticipant(tourneyID, opt.ParticipantID)
	return fmt.Sprintf(
		"Updated scores for %v in %v", participantData[0], tourneyData[0])
}
