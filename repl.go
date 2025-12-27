package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/cosmopolitics/cardingester/internal"
)

type Config struct {
	client          *http.Client
	cache           *cardingester.Cache
	selectedDataSet *string
}

type Command struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

func getCommands() map[string]Command {
	return map[string]Command{
		"help": {
			name:        "help",
			description: "prints commands and usage",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "(saves and) exits",
			callback:    commandExit,
		},
		"getdata": {
			name:        "getdata",
			description: "caches card data",
			callback:    commandGetCardData,
		},
		"search": {
			name:        "search",
			description: "searches for a card, prints info if one card is found, \n\tlists options if multiple are found",
			callback:    commandSearch,
		},
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func findScryfallBlob(url string, cache *cardingester.Cache, client *http.Client) ([]byte, int, error) {
	if blob, inDb := cache.Get(url); inDb {
		return blob, 0, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Set("User-Agent", "cardingest/0.0.1")
	request.Header.Set("Accept", "application/json")

	res, err := client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()
	blob, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("io error: %v", err)
	}

	cache.Add(url, blob)
	return blob, res.StatusCode, nil
}

func ScryfallUrlConstructor(params []string) (string, error) {
	baseUrl := "https://api.scryfall.com/cards/search?q="
	var url string = baseUrl

	for _, s := range params[:1] {
		if s == "search" {
			continue
		}
		if strings.Contains(s, ":") {
			p := strings.Split(s, ":")
			if len(p) < 2 {
				return "", fmt.Errorf("um split failed")
			}

			url = fmt.Sprintf("%s%s=%s", url, p[0], p[1])
		}
	}

	for _, s := range params[:1] {
		if s == "search" {
			continue
		}
		if strings.Contains(s, "order:") {
			p := strings.Split(url, "?")
			if len(p) < 2 {
				return "", fmt.Errorf("um split failed")
			}
			url = fmt.Sprintf("%s%s%s", p[0], "order=", p[1])
		}
	}

	for _, s := range params {
		if s == "search" {
			continue
		}
		url = fmt.Sprintf("%s%s", url, s)
	}

	return url, nil
}

func startRepl(cfg *Config) {
	reader := bufio.NewScanner(os.Stdin)

	for {
		color.RGB(203, 166, 247).Print("cardingester: ")
		reader.Scan()

		commands := cleanInput(reader.Text())
		if len(commands) == 0 {
			continue
		}

		cmd, exists := getCommands()[commands[0]]
		if exists {
			err := cmd.callback(cfg, commands)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
