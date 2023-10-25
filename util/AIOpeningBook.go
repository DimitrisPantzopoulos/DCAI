package util

import (
	"encoding/json"
	"fmt"
	"os"
)

type OpeningObject struct {
	Name         string
	OpeningMoves []string
}

type OpeningData struct {
	Openings map[string][][]string `json:"openings"`
}

func ReadOpeningJSONFile(filePath string) (*OpeningData, error) {
	var openingData OpeningData

	file, err := os.Open(filePath)
	if err != nil {
		return &openingData, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&openingData)
	if err != nil {
		return &openingData, err
	}

	return &openingData, nil
}

func AIOpeningBook(movesPlayed []string) (string, string) {
	filePath := "C:\\Users\\User\\Desktop\\Github\\Go\\util\\opening_data.json"

	openingData, err := ReadOpeningJSONFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return "Error", "Smth went wrong"
	}

	// FIND THE OPENING
	for openingName, moveSequences := range openingData.Openings {
		for _, moveSequence := range moveSequences {
			match := true
			for i := 0; i < len(movesPlayed) && i < len(moveSequence); i++ {
				if moveSequence[i] != movesPlayed[i] {
					match = false
					break // If there's a mismatch, break out of the loop
				}
			}
			if match {
				if len(moveSequence) > len(movesPlayed) {
					return openingName, moveSequence[len(movesPlayed)]
				} else {
					return openingName, "End of opening"
				}
			}
		}
	}

	return "Opening Not found :(", "Still not found :("
}
