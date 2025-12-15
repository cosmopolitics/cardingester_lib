package main

import (
	"bufio"
	"fmt"
	"os"
	"net/http"
)

type config struct {
	client *http.Client
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	cfg := &config{
		client: &http.Client{},
	}

	fmt.Println("Welcome to the Pokedex!")
	for {
		// Prompt
		green := "\033[32m"
		reset := "\033[0m"
		fmt.Print(green + "Pokedex: " + reset)

		reader.Scan()
		cleanText := cleanInput(reader.Text())
		if len(cleanText) == 0 {
			continue
		}

		// Do command
		commands := getCommands()
		if cmd, exist := commands[cleanText[0]]; exist {
			err := cmd.callback(cfg, cleanText)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("%s doesnt exist, 'help' for usage\n", cleanText[0])
		}
	}
}
