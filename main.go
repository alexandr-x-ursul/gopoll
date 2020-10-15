package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var polls map[string]Poll
var seqPoll int

// Poll structure stores info about a specific poll
type Poll struct {
	Question string         `json:"question"`
	Answers  map[string]int `json:"answers"`
}

// NewPoll is a structure used for posting new polls
type NewPoll struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

func getPolls(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(polls)
}

func getPoll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	poll, ok := polls[vars["id"]]
	if !ok {
		fmt.Fprintf(w, "No such poll found\n")
		return
	}
	json.NewEncoder(w).Encode(poll)
}

func nextPollID() string {
	seqPoll++
	return fmt.Sprint(seqPoll)
}

func createPoll(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	var newPoll NewPoll
	json.Unmarshal(body, &newPoll)
	resultPoll := Poll{
		Question: newPoll.Question,
		Answers:  make(map[string]int),
	}
	for _, v := range newPoll.Answers {
		resultPoll.Answers[v] = 0
	}
	polls[nextPollID()] = resultPoll
	fmt.Fprintf(w, "New Poll's ID: %d\n", seqPoll)
	json.NewEncoder(w).Encode(resultPoll)
}

func deletePoll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, ok := polls[id]; !ok {
		fmt.Fprintf(w, "No such poll found\n")
		return
	}
	delete(polls, id)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!\nThis application doesn't have a UI yet, only a REST API.\n")
}

func main() {
	polls = make(map[string]Poll)
	polls["11"] = Poll{
		Question: "why?",
		Answers: map[string]int{
			"idk":       0,
			"i do know": 1,
		},
	}
	polls["12"] = Poll{
		Question: "who is god?",
		Answers: map[string]int{
			"idk2":       1,
			"i do know2": 2,
		},
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/api/polls", getPolls).Methods("GET")
	router.HandleFunc("/api/poll", createPoll).Methods("POST")
	router.HandleFunc("/api/poll/{id}", deletePoll).Methods("DELETE")
	router.HandleFunc("/api/poll/{id}", getPoll).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
