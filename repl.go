package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cosmopolitics/cardingester/internal"
)

type Config struct {
	client *http.Client
	cache  *cardingester.Cache
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
	}
}

func commandHelp(cfg *Config, params []string) error {
	for _, cmd := range getCommands() {
		_, err := fmt.Printf("%s %s", cmd.name, cmd.description)
		if err != nil {
			return err
		}
	}
	return nil
}

func commandExit(cfg *Config, params []string) error {
	_, err := fmt.Printf("Closing the Pokedex... Goodbye!")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func startRepl(cfg *Config) {
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("cardingester > ")
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
