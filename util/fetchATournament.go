package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type TournamentWrapper struct {
	Data TournamentInfo `json:"data"`
}
type TournamentInfo struct {
	Attributes TournamentAttributes `json:"attributes"`
}
type TournamentAttributes struct {
	Name           string     `json:"name"`
	Game           string     `json:"game_name"`
	TournamentType string     `json:"tournament_type"`
	URL            string     `json:"full_challonge_url"`
	ImageURL       string     `json:"live_image_url"`
	StartAt        string     `json:"start_at"`
	RegOptions     RegOptions `json:"registration_options"`
}

type RegOptions struct {
	CheckIn int `json:"check_in_duration"`
}

func FetchATournament(tourneyID string) [7]string {
	var results [7]string
	//Request to the API
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v.json", tourneyID)
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
	var data TournamentWrapper
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	fmt.Println(data)
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		results[0] = ""
		return results
	}

	results[0] += data.Data.Attributes.Name
	results[1] += data.Data.Attributes.Game
	results[2] += data.Data.Attributes.TournamentType
	results[3] += data.Data.Attributes.URL
	results[4] += data.Data.Attributes.ImageURL
	results[5] += data.Data.Attributes.StartAt
	results[6] += strconv.Itoa(data.Data.Attributes.RegOptions.CheckIn)
	fmt.Println(results[6])
	return results
}
