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
	callback	func(*config) error
}

type config struct {
	nextLocationsURL *string
	prevLocationsURL *string
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
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cfg *config) error {
	commands := getCommands()
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
func commandMap(cfg *config) error {
	locations, err := pokeapi.FetchLocationAreas(cfg.nextLocationsURL)
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
func commandMapb(cfg *config) error {
	if cfg.prevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	locations, err := pokeapi.FetchLocationAreas(cfg.prevLocationsURL)
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
	}
}
