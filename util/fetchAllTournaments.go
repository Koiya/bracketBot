package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AllTournamentWrapper struct {
	Data []Info `json:"data"`
}
type Info struct {
	Attributes AllTournamentAttributes `json:"attributes"`
}
type AllTournamentAttributes struct {
	Name           string `json:"name"`
	Game           string `json:"game_name"`
	TournamentType string `json:"tournament_type"`
	URL            string `json:"full_challonge_url"`
}

func FetchAllTournaments() [3]string {
	var results [3]string
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
	var data AllTournamentWrapper
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
