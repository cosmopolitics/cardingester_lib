package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/cosmopolitics/cardingester/internal"
)

type Config struct {
	client *http.Client
	cache  *cardingester.Cache
	selectedDataSet *string
}

type Command struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands() map[string]Command {
	return map[string]Command{
		"help": {
			name:        "help",
			description: "- prints commands and usage",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "- (saves and) exits",
			callback:    commandExit,
		},
		"getdata": {
			name:        "getdata",
			description: "caches card data",
			callback:    commandGetCardData,
		},
	}
}

func commandHelp(cfg *Config, params []string) error {
	for _, cmd := range getCommands() {
		_, err := fmt.Printf("%s %s\n", cmd.name, cmd.description)
		if err != nil {
			return err
		}
	}
	return nil
}

func commandExit(cfg *Config, params []string) error {
  fmt.Println("Closing cardingester... Goodbye!")
	os.Exit(0)
	return nil
}

func commandGetCardData(cfg *Config, params []string) error {
	blob, err := findScryfallBlob("https://api.scryfall.com/bulk-data", cfg.cache, cfg.client)
	if err != nil {
		return err
	}
	var BOjson cardingester.Bulk_Option_Response
	err = json.Unmarshal(blob, &BOjson)
	if err != nil {
		return err
	}

	if params[1] == "help" {
		for _, entry := range BOjson.Data {
			fmt.Println(entry.Type)
			fmt.Println(entry.Description)
		}
		return nil
	}

	entryIndex := -1
	for i, entry := range BOjson.Data {
		if params[1] == entry.Type {
			entryIndex = i 
		}
	}
	if entryIndex == -1 {
		return fmt.Errorf("not a bulk card data option\n 'getdata help' for options")
	}


	cfg.selectedDataSet = &BOjson.Data[entryIndex].DownloadURI
	_, err = findScryfallBlob(*cfg.selectedDataSet, cfg.cache, cfg.client)


	return nil
}

func findScryfallBlob(url string, cache *cardingester.Cache, client *http.Client) ([]byte, error) {
	if blob, inDb := cache.Get(url); inDb {
		return blob, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Set("User-Agent", "cardingest/1.0")
	request.Header.Set("Accept", "application/json")

	res, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()
	blob, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io error: %v", err)
	}

	cache.Add(url, blob)
	return blob, nil
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

