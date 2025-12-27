package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmopolitics/cardingester/internal"
)

func commandHelp(cfg *Config, params []string) error {
	for _, cmd := range getCommands() {
		_, err := fmt.Printf("%s - %s\n", cmd.name, cmd.description)
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
	blob, code, err := findScryfallBlob("https://api.scryfall.com/bulk-data", cfg.cache, cfg.client)
	if code > 400 {

	}
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
	_, _, err = findScryfallBlob(*cfg.selectedDataSet, cfg.cache, cfg.client)

	return nil
}

func commandSearch(cfg *Config, params []string) error {
	url, err := ScryfallUrlConstructor(params)
	if err != nil {
		return err
	}

	blob, code, err := findScryfallBlob(url, cfg.cache, cfg.client)
	if err != nil {
		return err
	}

	if code < 400 {
		var search_json cardingester.Search_Response
		err = json.Unmarshal(blob, &search_json)
		if err != nil {
			return err
		}
		if search_json.Has_more == true {
			fmt.Println("there is a next page")
		}

		if len(search_json.Data) < 3 {
			for _, c := range search_json.Data {
				fmt.Println(c.Name)
				err := displayImage(c.ImageUris.Png, cfg)
				if err != nil {
					return err
				}
			}
		} else {
			for _, c := range search_json.Data {
				fmt.Println(c.Name)
			}
		}
	} else {
		var scryfall_error cardingester.Error_Response
		err := json.Unmarshal(blob, &scryfall_error)
		if err != nil {
			return nil
		}
		return fmt.Errorf("%s, %s", scryfall_error.Code, scryfall_error.Details)
	}

	return nil
}
