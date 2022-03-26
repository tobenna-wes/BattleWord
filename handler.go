package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	if err := http.ListenAndServe(":8083", nil); err != nil {
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

	// word := GuessWord()
	word := "world"

	guess := battleword.Guess{
		Guess: word,
		Shout: "Im fragmented",
	}

	err = json.NewEncoder(w).Encode(guess)
	if err != nil {
		log.Println(err)
		return
	}
	prevGuessesJSON, _ := json.Marshal(prevGuesses)
	fmt.Printf("Making random guess for game %s, turn %d: %s\n", r.Header.Get(battleword.GuessIDHeader), len(prevGuesses.GuessResults), word)
	fmt.Printf("Request ID %s. Body: %s\n", r.Header.Get(battleword.GuessIDHeader), prevGuessesJSON)
}

func HandlePing(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Handle Result Triggered")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
	}

	def := &battleword.PlayerDefinition{
		Name:                "Lucid Thinker",
		Description:         "Quick",
		ConcurrentConnLimit: 10,
		Colour:              "#006179",
	}

	err := json.NewEncoder(w).Encode(def)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func HandleResult(w http.ResponseWriter, r *http.Request) {
	var finalState battleword.PlayerMatchResults
	err := json.NewDecoder(r.Body).Decode(&finalState)
	if err != nil {
		fmt.Println(err)
		return
	}

	var us battleword.Player
	found := false
	for _, player := range finalState.Results.Players {
		if player.ID == finalState.PlayerID {
			us = player
			found = true
		}
	}

	if !found {
		log.Println("We weren't in the results. strange")
		return
	}

	gamesWon := 0
	for _, game := range us.GamesPlayed {
		if game.Correct {
			gamesWon++
		}
	}

	// finalStateJSON, _ := json.Marshal(finalState)

	fmt.Printf("Match %s concluded, we got %d words right", finalState.Results.UUID, gamesWon)

	// log.Printf("Match %s concluded, we got %d words right. Body: %s", finalState.Results.UUID, gamesWon, finalStateJSON)
}

func GuessWord() string {

	return battleword.CommonWords[rand.Intn(len(battleword.CommonWords))]
}
