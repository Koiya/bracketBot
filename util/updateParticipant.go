package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

/* '{"data":{"type":"Participants","attributes":{"name":"As","seed":1,"misc":"","email":"","username":""}}}' */
func UpdateParticipant(tourneyID, participantID string, opt Options) string {

	//Request to the API
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/participants/%v.json", tourneyID, participantID)
	var updatedAttributes string
	if opt.Seed == 0 {
		opt.Seed = 1
	}
	if opt.Misc != "" {
		updatedAttributes += `"misc" : "` + opt.Misc + `",`
	}
	if opt.Name != "" {
		updatedAttributes += `"name" : "` + opt.Name + `",`
	}
	if opt.Username != "" {
		updatedAttributes += `"username" : "` + opt.Username + `",`
	}
	if opt.Email != "" {
		updatedAttributes += `"email" : "` + opt.Email + `",`
	}
	fmt.Println(opt)
	fmt.Println(updatedAttributes)
	requestBody := fmt.Sprintf(`{
		"data": {
		"type" : "Participants",
		"attributes" : {
				%v
				"seed": "%d"
			}
		}
	}`, updatedAttributes, opt.Seed)
	fmt.Println(requestBody)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Body:", string(body))

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return "Name already exists"
	}
	tourneyData := FetchATournament(tourneyID)
	return fmt.Sprintf(
		"Updated information of %v in %v", opt.Name, tourneyData[0],
	)
}
