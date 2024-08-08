package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type ParticipantWrapper struct {
	Data ParticipantInfo `json:"data"`
}
type ParticipantInfo struct {
	Id         string                `json:"id"`
	Attributes ParticipantAttributes `json:"attributes"`
}
type ParticipantAttributes struct {
	Name     string `json:"name"`
	Seed     int    `json:"seed"`
	Misc     string `json:"misc"`
	Username string `json:"username"`
}

func FetchParticipant(tourneyID, participantID string) [4]string {
	var results [4]string
	//Request to the API
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/participants/%v.json", tourneyID, participantID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
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
	var data ParticipantWrapper
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error", err.Error())
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		results[0] = "Error getting data"
		return results
	}
	results[0] = data.Data.Attributes.Name
	results[1] = data.Data.Id
	results[2] = strconv.Itoa(data.Data.Attributes.Seed)
	//fmt.Println(data.Data)

	return results
}
