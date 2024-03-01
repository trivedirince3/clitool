package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"flag"
	"os"
)

type Definition struct {
	Shortdef []string `json:"shortdef"`
	Hwi      struct {
		Hw string `json:"hw"`
		Ipa string `json:"omitempty"`
	} `json:"hwi"`
}

func wordLookup(word string, apiKey string) (string, error) {
	url := fmt.Sprintf("https://www.dictionaryapi.com/api/v3/collegiate/definitions?key=%s&word=%s", apiKey, word)
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get response: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", response.StatusCode)
	}

	var data []Definition
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(data) == 0 {
		return "Word not found.", nil
	}

	definition := data[0].Shortdef[0]
	pronunciation := data[0].Hwi.Ipa
	if pronunciation == "" {
		pronunciation = "(no pronunciation available)"
	}
	return fmt.Sprintf("%s (%s): %s", pronunciation, data[0].Hwi.Hw, definition), nil
}

func main() {
	word := flag.String("word", "", "The word to look up.")
	flag.Parse()

	if *word == "" {
		fmt.Println("Missing required argument: word")
		os.Exit(1)
	}

	apiKey := os.Getenv("MW_API_KEY") // Read API key from environment variable
	if apiKey == "" {
		fmt.Println("Missing environment variable: MW_API_KEY (Merriam-Webster API key)")
		os.Exit(1)
	}

	definition, err := wordLookup(*word, apiKey)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(definition)
}
