package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"bytes"
	"github.com/dolmen-go/kittyimg"
	"github.com/fatih/color"
	"github.com/sunshineplan/imgconv"
	"image"

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
	baseUrl := "https://api.scryfall.com/cards/search?order=name&q="
	var url string = baseUrl

	for _, word := range params[1:] {
		if strings.Contains(word, "order:") {
			sortmode := strings.Split(word, "order:")
			preurl := strings.Split(url, "order=")
			queryurl := strings.Split(url, "q=")

			url = fmt.Sprintf("%sorder=%s&q=%s", preurl[0], sortmode[1], queryurl[1])

		} else {
			url = fmt.Sprintf("%s%s", url, word)
		}
	}

	fmt.Println(url)
	return url, nil
}

func displayImage(url string, cfg *Config) error {
	img, _, err := findScryfallBlob(url, cfg.cache, cfg.client)
	if err != nil {
		return err
	}
	decodedImage, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return err
	}

	imageBuffer := new(bytes.Buffer)
	err = imgconv.Write(imageBuffer, decodedImage, &imgconv.FormatOption{Format: imgconv.PNG})
	if err != nil {
		return err
	}

	err = kittyimg.Transcode(os.Stdout, imageBuffer)
	fmt.Print("\n")
	if err != nil {
		return err
	}
	return nil
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
