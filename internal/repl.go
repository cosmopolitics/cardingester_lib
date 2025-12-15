package cardingester

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Command struct {
	name string
	description string
	callback func(*config, []string) error
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands() map[string]Command {
	return map[string]Command{
		"help": {
			name: "help",
			description: "- prints commands and usage",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "- (saves and) exits",
			callback: commandExit,
		},
	}
}

func commandHelp(cfg *config, params []string) error {
	for _, cmd := range getCommands() {
		_, err := fmt.Printf("%s %s", cmd.name, cmd.description)
		if err != nil {
			return err
		}
	}
	return nil
}

func commandExit(cfg *config, params []string) error {
	_, err := fmt.Printf("Closing the Pokedex... Goodbye!")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func (cfg *config) getBlob(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// setting a header on the new request
	request.Header.Set("User-Agent", "cardingest/1.0")
	request.Header.Set("Accept", "application/json")

	res, err := cfg.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()
	blob, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io error: %v", err)
	}

	return blob, nil
}
