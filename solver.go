package main

import (
	"fmt"
	"strings"

	"github.com/brensch/battleword"
)

var vowels = []string{"a", "e", "i", "o", "u"}

var commonLetters = []string{"t", "n", "l", "c", "d", "s"}

// var vowels = []string{}
// var commonLetters = []string{}
var allLetters = "abcdefghijklmnopqrstuvqrstuvwxyz"

type letterStrong struct {
	Letter  string    `json:"letter,omitempty"`
	Weight  []float32 `json:"weight,omitempty"`
	Locaion []float32 `json:"locaion,omitempty"`
}

type WordStrong struct {
	word   string  `json:"word,omitempty"`
	Weight float32 `json:"weight,omitempty"`
}

type letterStringList struct {
	List  []letterStrong `json:"list,omitempty"`
	twins []string
}

func removeDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func hasDuplicateValues(intSlice []string) (string, int) {
	keys := make(map[string]bool)

	for i, entry := range intSlice {
		_, value := keys[entry]

		if !value {
			keys[entry] = true

		} else {
			return entry, i
		}
	}
	return "", 0
}

func GetDefaultLetterWeight(letter string) []float32 {

	if len(letter) != 1 {
		return []float32{1, 1, 1, 1, 1}
	}

	for _, vowel := range vowels {
		if vowel == letter {
			return []float32{3, 3, 3, 3, 3}
		}
	}

	for _, common := range commonLetters {
		if common == letter {
			return []float32{2, 2, 2, 2, 2}
		}
	}

	return []float32{1, 1, 1, 1, 1}

}

func GetDefaultLetterWeightFromList() []letterStrong {

	commonWords := battleword.CommonWords

	keys := make(map[string]float32)

	max := []float32{0, 0, 0, 0, 0}

	location := 1
	for _, entry := range commonWords {

		wordSingles := strings.Split(entry, "")
		_, value := keys[wordSingles[location]]

		if !value {
			keys[wordSingles[location]] = 0

		} else {
			keys[wordSingles[location]] += 1
		}

		val, _ := keys[wordSingles[location]]

		if val > max[location] {
			max[location] = val
		}
	}

	for key := range keys {
		keys[key] = keys[key] / max[location]
	}

	fmt.Println("keys", keys, len(commonWords), max[location])
	return []letterStrong{}

}

func CreateAllDefaultWiefghts() letterStringList {

	wordSingles := strings.Split(allLetters, "")
	allStringList := letterStringList{}

	for _, letter := range wordSingles {
		allStringList.List = append(allStringList.List,
			letterStrong{
				Letter:  letter,
				Weight:  GetDefaultLetterWeight(letter),
				Locaion: []float32{1, 1, 1, 1, 1},
			})
	}

	// fmt.Println(len(allStringList.List))

	return allStringList
}

func SetingleWiefghts(lists letterStringList, new letterStrong) letterStringList {

	newList := letterStringList{}
	for _, list := range lists.List {

		if list.Letter == new.Letter {
			newList.List = append(newList.List, new)
		} else {
			newList.List = append(newList.List, list)
		}

	}

	return newList
}

func GetSingleWiefghts(lists letterStringList, letter string) letterStrong {

	for _, list := range lists.List {

		if list.Letter == letter {
			return list
		}

	}

	return letterStrong{
		Letter: letter,
		Weight: []float32{1, 1, 1, 1, 1},
	}
}

func UpdateDefaultWiefghts(list letterStringList, w *battleword.PlayerGameState) letterStringList {

	if w == nil {
		return list
	}

	games := w.GuessResults

	newlist := list

	for _, game := range games {
		wordSingles := strings.Split(game.Guess, "")
		fmt.Println("last game", game.Result, game.Guess)

		twins, twinLocation := hasDuplicateValues(wordSingles)
		for i, letter := range wordSingles {
			lastStrong := GetSingleWiefghts(newlist, letter)

			// fmt.Println("twins:", twins, ":", letter, twinLocation, i)
			// fmt.Println("old letter update", lastStrong)

			if game.Result[i] == 0 {

				lastStrong.Locaion[i] = 0
				lastStrong.Weight[i] = 0

				if twins == letter && twinLocation == i {

				} else {
					for j, loc := range lastStrong.Locaion {
						if loc != 2 {
							lastStrong.Locaion[j] = 0
							lastStrong.Weight[j] = 0
						}
					}
				}

			}

			if game.Result[i] == 1 {

				lastStrong.Locaion[i] = 0
				lastStrong.Weight[i] = 0
				for j, loc := range lastStrong.Locaion {
					if loc != 0 {
						lastStrong.Locaion[j] = 1
						lastStrong.Weight[j] = 40
					}
				}

				// lastStrong.Weight = []int{20, 20, 20, 20, 20}
				// lastStrong.Locaion = []int{1, 1, 1, 1, 1}

			}

			if game.Result[i] == 2 {

				if twins == letter && twinLocation == i {

				} else {
					lastStrong.Weight = []float32{0, 0, 0, 0, 0}
				}

				lastStrong.Weight[i] = 60
				lastStrong.Locaion[i] = 2
			}

			fmt.Println("new letter update", lastStrong)
			newlist = SetingleWiefghts(newlist, lastStrong)
		}

	}

	fmt.Println("--------------------------")
	for _, result := range newlist.List {
		fmt.Println(result)
	}
	fmt.Println("--------------------------")

	return newlist
}

func Solve(w battleword.PlayerGameState) string {

	fmt.Println(w.GuessResults)

	for _, result := range w.GuessResults {
		fmt.Println(result.Guess)
	}

	fmt.Println(len(w.GuessResults) + 1)

	// if len(w.GuessResults) == 0 {
	// 	return CreateGuess(&w).Letter
	// }
	return CreateGuess(&w).word
}

func GetLetterWeight(letter string, wig letterStringList) letterStrong {

	for _, listed := range wig.List {
		if listed.Letter == letter {
			return listed
		}
	}

	return letterStrong{
		Letter: letter,
		Weight: []float32{1, 1, 1, 1, 1},
	}
}

func GetWordWeight(word string, wig letterStringList) WordStrong {
	wordSingles := strings.Split(word, "")

	// wordSingles = removeDuplicateValues(wordSingles)

	totalW := float32(0)
	for i, single := range wordSingles {
		// fmt.Println(single, GetLetterWeight(single))
		letterWight := GetLetterWeight(single, wig)

		if letterWight.Locaion[i] == 2 {
			totalW += letterWight.Weight[i] + 40

		} else if letterWight.Locaion[i] == 1 {
			totalW += letterWight.Weight[i] + 30

		} else {
			if letterWight.Weight[i] != 0 {
				totalW += letterWight.Weight[i]
			}
		}

	}

	return WordStrong{
		word:   word,
		Weight: totalW,
	}
}

func CreateGuess(w *battleword.PlayerGameState) WordStrong {

	commonWords := battleword.CommonWords
	// commonWords := append(battleword.AllWords, battleword.CommonWords...)
	weights := CreateAllDefaultWiefghts()
	weights = UpdateDefaultWiefghts(weights, w)

	// fmt.Println()

	// maxWordWight := 0
	maxWord := WordStrong{
		word:   "",
		Weight: 0,
	}

	for _, commonWord := range commonWords {

		w := GetWordWeight(commonWord, weights)

		if w.Weight >= maxWord.Weight {
			maxWord.Weight = w.Weight
			maxWord.word = commonWord
			// fmt.Println(commonWord, w)
		}

	}

	return maxWord
}
