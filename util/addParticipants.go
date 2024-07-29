package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type PostAddParticipantWrapper struct {
	Data []PostAddParticipant `json:"data"`
}

type PostAddParticipant struct {
	Type string `json:"type"`
}

type PostAddParticipantAttributes struct {
	Name string `json:"name"`
	Seed string `json:"seed"`
	//Optional
	Misc     string `json:"misc"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type options struct {
	seed     int
	misc     string
	email    string
	username string
}
type option func(*options)

/* '{"data":{"type":"Participants","attributes":{"name":"As","seed":1,"misc":"","email":"","username":""}}}' */
func AddParticipants(tourneyID, name string, etc ...option) {
	var seed int
	var misc string
	var email string
	var username string

	fmt.Println("Tourney ID: " + tourneyID)
	fmt.Println("Name: " + name)
	//Request to the API
	if tourneyID == "" {
		return
	}
	var URL string
	URL = fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%v/participants.json", tourneyID)
	o := &options{}
	if o.seed > 0 {
		seed = o.seed
	} else {
		seed = 1
	}

	fmt.Println("URL" + URL)
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
	}`, name, seed, misc, email, username)
	//reqBody := `{"data":{"type":"Participants","attributes":{"name":"As","seed":1,"misc":"","email":"","username":""}}}`
	fmt.Println(requestBody)

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

	if resp.StatusCode != http.StatusCreated {

		fmt.Println(resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("response Body:", string(body))

}
