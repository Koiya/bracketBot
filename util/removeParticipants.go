package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func RemoveParticipants(tourneyID, name string) string {

	//Request to the API
	if tourneyID == "" {
		return "No ID inputted"
	}

	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/participants/%v.json", tourneyID, name)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", URL, nil)
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
	if resp.StatusCode != 204 || resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(string(body))
		return "Error trying to remove a participant"
	}

	return fmt.Sprintf(
		"Removed %v from %v", name, tourneyID,
	)
}
