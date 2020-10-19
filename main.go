package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"sync"
)

var polls map[string]Poll
var seqPoll int
var mutex sync.Mutex

// Poll structure stores info about a specific poll
type Poll struct {
	Question string          `json:"question"`
	Answers  map[string]int  `json:"answers"`
	Answered map[string]bool `json:"-"`
}

// NewPoll is a structure used for posting new polls
type NewPoll struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

// Vote is a structure used for posting a vote on a poll
type Vote struct {
	Answer string `json:"answer"`
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
	mutex.Lock()
	seqPoll++
	mutex.Unlock()
	return fmt.Sprint(seqPoll)
}

func createPoll(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var newPoll NewPoll
	if err := decoder.Decode(&newPoll); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resultPoll := Poll{
		Question: newPoll.Question,
		Answers:  make(map[string]int),
		Answered: make(map[string]bool),
	}
	for _, v := range newPoll.Answers {
		resultPoll.Answers[v] = 0
	}
	id := nextPollID()
	polls[id] = resultPoll
	fmt.Fprintf(w, "New Poll's ID: %s\n", id)
	json.NewEncoder(w).Encode(resultPoll)
}

func deletePoll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := polls[id]; !ok {
		http.Error(w, "No such poll found.", 400)
		return
	}
	delete(polls, id)
}

func vote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mutex.Lock()
	defer mutex.Unlock()
	var vote Vote
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&vote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	poll, ok := polls[vars["id"]]
	if !ok {
		http.Error(w, "No such poll found.", 400)
		return
	}
	if _, ok := poll.Answers[vote.Answer]; !ok {
		fmt.Fprintf(w, "No such answer found.\n")
		return
	}
	ip := userIP(r)
	if _, exists := poll.Answered[ip]; exists {
		http.Error(w, "You've already voted!", http.StatusBadRequest)
		return
	}
	poll.Answered[ip] = true
	poll.Answers[vote.Answer]++
}

func userIP(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	ipAddress = ipAddress[:strings.LastIndex(ipAddress, ":")]
	return ipAddress
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
		Answered: map[string]bool{},
	}
	polls["12"] = Poll{
		Question: "who is god?",
		Answers: map[string]int{
			"idk2":       1,
			"i do know2": 2,
		},
		Answered: map[string]bool{},
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/api/polls", getPolls).Methods("GET")
	router.HandleFunc("/api/poll", createPoll).Methods("POST")
	router.HandleFunc("/api/poll/{id}", deletePoll).Methods("DELETE")
	router.HandleFunc("/api/poll/{id}", getPoll).Methods("GET")
	router.HandleFunc("/api/poll/{id}", vote).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
