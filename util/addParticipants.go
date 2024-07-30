package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

//type PostAddParticipantWrapper struct {
//	Data []PostAddParticipant `json:"data"`
//}
//
//type PostAddParticipant struct {
//	Type string `json:"type"`
//}
//
//type PostAddParticipantAttributes struct {
//	Name string `json:"name"`
//	Seed string `json:"seed"`
//	//Optional
//	Misc     string `json:"misc"`
//	Email    string `json:"email"`
//	Username string `json:"username"`
//}

type Options struct {
	Seed     int
	Misc     string
	Email    string
	Username string
}

/* '{"data":{"type":"Participants","attributes":{"name":"As","seed":1,"misc":"","email":"","username":""}}}' */
func AddParticipants(tourneyID, name string, opt Options) string {
	fmt.Println("Tourney ID: " + tourneyID)
	fmt.Println("Name: " + name)
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
				"email" : "%v",
				"username" : "%v"
			}
		}
	}`, name, opt.Seed, opt.Misc, opt.Email, opt.Username)

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

	if resp.StatusCode != http.StatusCreated {
		fmt.Println(resp.StatusCode)
		return "Name already exists"
	}

	return fmt.Sprintf(
		"Added %v to %v", name, tourneyID,
	)
}
