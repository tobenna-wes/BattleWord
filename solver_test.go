package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/brensch/battleword"
	"github.com/sirupsen/logrus"
)

func TestGetLetterWeight(t *testing.T) {
	weights := CreateAllDefaultWiefghts()
	w := GetLetterWeight("q", weights)
	fmt.Println("w:", w)
}

func TestGetWordWeight(t *testing.T) {
	weights := CreateAllDefaultWiefghts()
	w := GetWordWeight("aa", weights)
	fmt.Println("w:", w)
}

func TestCreateGuess(t *testing.T) {
	w := CreateGuess(nil)
	fmt.Println("guess:", w)
}

func TestCreateGame(t *testing.T) {
	ServerStart()
}

func TestCreateEgine(t *testing.T) {
	RunEnigneServer("http://localhost:8083/", 5, 5, 1)
}

func RunBattleServerS() {

}

func RunEnigneServer(PlayerURIsJoined string, NumLetters, NumRounds, NumGames int) {

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	log.Info("started game")
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	playerURIs := strings.Split(PlayerURIsJoined, ",")

	if playerURIs[0] == "" {
		log.Error("you need to define player endpoints")
		return
	}

	// windows can't contain : in filenames. stitchup
	filename := fmt.Sprintf("results-%s.json", time.Now().Format("20060102-150405-0700"))
	f, err := os.Create(filename)
	if err != nil {
		log.WithError(err).WithField("file", filename).Error("couldn't create file")
		return
	}
	defer f.Close()

	match, err := battleword.InitMatch(log, battleword.AllWords, battleword.CommonWords, playerURIs, NumLetters, NumRounds, NumGames)
	if err != nil {
		log.WithError(err).Error("got error initialising game")
		return
	}

	match.Start(context.Background())
	match.Broadcast()

	err = json.NewEncoder(f).Encode(match.Snapshot())
	if err != nil {
		log.WithError(err).Error("couldn't write to file")
		return
	}

	log.WithField("file", filename).Println("final result saved to file")
}
