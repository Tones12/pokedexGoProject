package main

import (
	"strings"
	"fmt"
	"os"
	"bufio"
	"net/http"
	"encoding/json"
	"io"
	"github.com/tones12/pokedexgoproject/internal/pokeapi"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
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
		description: "Displays Pokemon locations, 20 at a time",
		callback:	 commandMap,
	},
	}
}

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*config) error {
	commands := getCommands()
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(*config) error {
	
	pokeapi.FetchLocationAreas()
	return nil
}

func cleanInput(text string) []string {
	var words []string
	lowerText := strings.ToLower(text)
	words = strings.Fields(lowerText)
	return words
}

func startRepl() {
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
			err := command.callback()
			if err != nil {
				fmt.Println("Error executing command: ", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
