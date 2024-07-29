package cmd

type Participant struct {
	Name     string
	Seed     int
	Misc     string
	Email    string
	Username string
}

func init() {

}
func addParticipants(TourneyID, seed int, name, misc, email, username string) {
	/*	url := fmt.Sprintf("https://api.challonge.com/v2.1/tournaments/%d/participants.json", TourneyID)
		d := Participant{
			name,
			seed,
			misc,
			email,
			username,
		}
		b, err := json.Marshal(d)*/
}

//https://api.challonge.com/v2.1/tournaments/{tourneyid}/participants.json
