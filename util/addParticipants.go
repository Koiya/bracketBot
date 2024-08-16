package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Options struct {
	Name     string
	Seed     int
	Misc     string
	Email    string
	Username string
}

func AddParticipants(tourneyID string, opt Options) string {

	//Request to the API
	if tourneyID == "" {
		return "No ID inputted"
	}
	if opt.Seed == 0 {
		opt.Seed = 1
	}
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/participants.json", tourneyID)
	requestBody := fmt.Sprintf(`{
		"data": {
		"type" : "Participants",
		"attributes" : {
				"name" : "%v",
				"seed" : %d,
				"misc" : "%v",
				"email" : "",
				"username" : ""
			}
		}
	}`, opt.Name, opt.Seed, opt.Misc)

	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", GetTOML("API.token"))

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Body:", string(body))

	if resp.StatusCode == 422 {
		fmt.Println(resp.StatusCode)
		return "Participant already added."
	} else if resp.StatusCode != http.StatusCreated {
		fmt.Println(resp.StatusCode)
		return "Error adding in participant. Please check parameters."
	}
	tourneyData := FetchATournament(tourneyID)
	return fmt.Sprintf(
		"Successfully added %v to %v", opt.Name, tourneyData[0],
	)
}
