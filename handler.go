package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/brensch/battleword"
)

func ServerStart() {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("worr")

	})

	http.HandleFunc("/guess", HandleGuess)
	http.HandleFunc("/results", HandleResult)
	http.HandleFunc("/ping", HandlePing)

	fmt.Printf("Starting server at port\n")
	if err := http.ListenAndServe(":8056", nil); err != nil {
		log.Fatal(err)
	}
}

func HandleGuess(w http.ResponseWriter, r *http.Request) {
	// Double check it's a post request being made
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	var prevGuesses battleword.PlayerGameState

	err := json.NewDecoder(r.Body).Decode(&prevGuesses)
	if err != nil {
		log.Println(err)
		return
	}

	word := ""
	guess := battleword.Guess{
		Guess: word,
		Shout: "Im fragmented",
	}

	fmt.Fprintf(w, "We started")

	err = json.NewEncoder(w).Encode(guess)
	if err != nil {
		log.Println(err)
		return
	}
	prevGuessesJSON, _ := json.Marshal(prevGuesses)
	log.Printf("Making random guess for game %s, turn %d: %s\n", r.Header.Get(battleword.GuessIDHeader), len(prevGuesses.GuessResults), word)
	log.Printf("Request ID %s. Body: %s\n", r.Header.Get(battleword.GuessIDHeader), prevGuessesJSON)
}

func HandleResult(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Handle Result Triggered")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
	}

	def := &battleword.PlayerDefinition{
		Name:                "solvo",
		Description:         "the magnificent",
		ConcurrentConnLimit: 10,
		Colour:              "#596028",
	}

	err := json.NewEncoder(w).Encode(def)
	if err != nil {
		log.Println(err)
		return
	}
}

func HandlePing(w http.ResponseWriter, r *http.Request) {

}
