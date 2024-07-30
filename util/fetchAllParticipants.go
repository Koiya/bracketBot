package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type AllParticipantWrapper struct {
	Data []AllParticipantInfo `json:"data"`
}
type AllParticipantInfo struct {
	Id         string                   `json:"id"`
	Attributes AllParticipantAttributes `json:"attributes"`
}
type AllParticipantAttributes struct {
	Name     string `json:"name"`
	Seed     int    `json:"seed"`
	Misc     string `json:"misc"`
	Username string `json:"username"`
}

func FetchAllParticipants(tourneyID string) [4]string {
	var results [4]string
	//Request to the API
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/participants.json?page=1&per_page=25", tourneyID)
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
	var data AllParticipantWrapper
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error", err.Error())
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		results[0] = "Error getting data"
		return results
	}
	for _, v := range data.Data {
		results[0] += fmt.Sprintf("%v (%v)\n\n\n", v.Attributes.Name, v.Id)
		results[1] += v.Attributes.Misc + "\n\n\n"
		results[2] += strconv.Itoa(v.Attributes.Seed) + "\n\n\n"
	}
	//fmt.Println(data.Data)

	return results
}
