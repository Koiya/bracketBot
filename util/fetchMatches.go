package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type AllMatchesWrapper struct {
	Data []AllMatchesInfo `json:"data"`
}
type AllMatchesInfo struct {
	Id         string               `json:"id"`
	Attributes AllMatchesAttributes `json:"attributes"`
}
type AllMatchesAttributes struct {
	State               string                `json:"state"`
	Round               int                   `json:"round"`
	Scores              string                `json:"scores"`
	PointsByParticipant []PointsByParticipant `json:"points_by_participant"`
}

type PointsByParticipant struct {
	ParticipantID int `json:"participant_id"`
}

func FetchAllMatches(tourneyID string) [4]string {
	var results [4]string
	var participantsName [2]int
	//Request to the API
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/matches.json", tourneyID)
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
	var data AllMatchesWrapper
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
		for i, v := range v.Attributes.PointsByParticipant {
			participantsName[i] = v.ParticipantID
		}
		results[0] += strconv.Itoa(participantsName[0]) + "\n\n\n"
		results[1] += v.Attributes.Scores + "\n\n\n"
		results[2] += strconv.Itoa(participantsName[1]) + "\n\n\n"
	}
	fmt.Println(data.Data)

	return results
}
