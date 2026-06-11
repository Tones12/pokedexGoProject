package main

import (
	"strings"
	"fmt"
	"os"
	"bufio"
	"github.com/tones12/pokedexgoproject/internal/pokeapi"
)

type cliCommand struct {
	name		string
	description	string
	callback	func(*config, string) error
}

type config struct {
	nextLocationsURL *string
	prevLocationsURL *string
	pokeapiClient    *pokeapi.Client
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
	"exit": {
		name:		 "exit",
		description: "Exit the Pokedex",
		callback:	 commandExit,
	},
	"help": {
		name:		 "help",
		description: "Displays a help message",
		callback:	 commandHelp,
	},
	"map": {
		name:		 "map",
		description: "Displays Pokemon locations, 20 at a time, and the next 20 locations",
		callback:	 commandMap,
	},
	"mapb": {
		name:		 "mapb",
		description: "Displays the previous 20 locations",
		callback:	 commandMapb,
	},
	"explore": {
		name:		 "explore",
		description: "Explore + area name returns pokemon found in that area",
		callback:	 commandExplore,
	},
	}
}

func commandExit(cfg *config, name string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cfg *config, name string) error {
	commands := getCommands()
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
func commandMap(cfg *config, name string) error {
	locations, err := cfg.pokeapiClient.FetchLocationAreas(cfg.nextLocationsURL)
	if err != nil {
		return err
	}
	cfg.nextLocationsURL = locations.Next
	cfg.prevLocationsURL = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}
func commandMapb(cfg *config, name string) error {
	if cfg.prevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	locations, err := cfg.pokeapiClient.FetchLocationAreas(cfg.prevLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locations.Next
	cfg.prevLocationsURL = locations.Previous
	
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}
func commandExplore(cfg *config, name string) error {
	if cfg.prevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	locations, err := cfg.pokeapiClient.FetchLocationAreas(cfg.prevLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locations.Next
	cfg.prevLocationsURL = locations.Previous
	
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}
func cleanInput(text string) []string {
	var words []string
	lowerText := strings.ToLower(text)
	words = strings.Fields(lowerText)
	return words
}
func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := cleanInput(scanner.Text())
		if len(userInput) == 0 {
			continue
		}
		if len(userInput) == 1 {
			commandName := userInput[0]
			commands := getCommands()
			if command, ok := commands[commandName]; ok {
				err := command.callback(cfg)
			if err != nil {
				fmt.Println("Error executing command: ", err)
			}
			} else {
			fmt.Println("Unknown command")
			}
		} else if len(userInput) == 2 {
			commandName := userInput[0]
			locName := userInput[1]
			commands := getCommands()
			if command, ok := commands[commandName]; ok {
				err := command.callback(cfg, locName)
				if err != nil {
					fmt.Println("Error executing command: ", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
