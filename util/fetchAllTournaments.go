package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DataWrapper struct {
	Data []Info `json:"data"`
}
type Info struct {
	Attributes Attributes `json:"attributes"`
}
type Attributes struct {
	Name           string `json:"name"`
	Game           string `json:"game_name"`
	TournamentType string `json:"tournament_type"`
	URL            string `json:"full_challonge_url"`
}

func FetchAllTournaments() [4]string {
	/*
		data['tournament_type']
		data['full_challonge_url']
		data['game_name']

	*/
	var results [4]string
	//Request to the API
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.challonge.com/v2.1/tournaments.json?page=1&per_page=25", nil)
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
	var data DataWrapper
	//json.Unmarshal(body, &response)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	for _, v := range data.Data {
		results[0] += fmt.Sprintf("[%v](%v)\n\n\n", v.Attributes.Name, v.Attributes.URL)
		results[1] += v.Attributes.Game + "\n\n\n"
		results[2] += v.Attributes.TournamentType + "\n\n\n"
	}
	fmt.Println(data.Data)
	//fmt.Println(data.Result.Att.Game)
	//fmt.Println("Tournament Type:", , "Game: ", data.Game)
	return results
}
