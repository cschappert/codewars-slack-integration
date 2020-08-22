package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
)

var api *slack.Client
var channelId string
var codewarsIconUrl string = "https://www.codewars.com/assets/logos/logo-square-paper-bg-c3d2b1eb4fb35d75b0c0c0e3b74616fab527afdce9d1d3184624cf0b4e950357.jpg"

type ChallengeDetail struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CodeChallenge struct {
	Id          string `json:"id"`
	CreatedById string `json:"created_by_id"`
}

type Solution struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type Event struct {
	Action        string         `json:"action"`
	CodeChallenge *CodeChallenge `json:"code_challenge,omitempty"`
	Solution      *Solution      `json:"solution,omitempty"`
}

func getChallengeDetail(challenge string) *ChallengeDetail {
	url := fmt.Sprintf("https://www.codewars.com/api/v1/code-challenges/%s", challenge)
	res, _ := http.Get(url)
	resBody, _ := ioutil.ReadAll(res.Body)

	var challengeDetail ChallengeDetail
	json.Unmarshal(resBody, &challengeDetail)

	return &challengeDetail
}

func postMessage(user, challenge string) {
	cd := getChallengeDetail(challenge)
	msg := fmt.Sprintf("User *%s* completed coding challenge *%s*\nTry it! %s", user, cd.Name, cd.Url)
	_, _, _, err := api.SendMessage(
		channelId,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionUsername("Codewars"),
		slack.MsgOptionIconURL(codewarsIconUrl),
	)

	if err != nil {
		fmt.Println(err)
	}
}

// Return 200
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK!")
}

// Handle Codewars events
func CodewarsHookHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var event Event
	json.Unmarshal(reqBody, &event)

	if event.Action == "solution_finalized" {
		fmt.Println("enter")
		postMessage(event.Solution.UserId, event.CodeChallenge.Id)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
	// create slack API client
	api = slack.New(os.Getenv("SLACK_TOKEN"))
	channelId = os.Getenv("CHANNEL_ID")

	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/hook/codewars", CodewarsHookHandler)

	log.Fatal(http.ListenAndServe(":8000", r))
}
