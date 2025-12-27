package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"net/url"

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
	baseUrl := "https://api.scryfall.com/cards/"
	query := url.Values{}
	finalUrl := baseUrl + "search?" 

	var filter []string
	var name_text []string
	for _, word := range params[1:] {
		if strings.Contains(word, ":") {
			filter = append(filter, word)

		} else if strings.Contains(word, "order:") {
			sort := strings.Split(word, "order:")
			finalUrl = finalUrl + "order%3D" + sort[1] + "&"

		} else {
			name_text = append(name_text, word)
		}
	}

	plain_text := strings.Join(name_text, " ")
	query.Add("q", plain_text)

	for _, f := range filter {
		fparts := strings.Split(f, ":")
		if len(fparts) < 2 {
			return "", fmt.Errorf("filter split failed %v", fparts)
		}
		query.Add(fparts[0], fparts[1])
	}

	finalUrl = finalUrl + "&" + query.Encode()
	fmt.Println(finalUrl)

	return finalUrl, nil
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
	if err != nil {
		return err
	}
	fmt.Print("\n")
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
