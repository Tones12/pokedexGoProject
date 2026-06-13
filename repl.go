package main

import (
	"strings"
	"fmt"
	"os"
	"bufio"
	"github.com/tones12/pokedexgoproject/internal/pokeapi"
	"math/rand"
)

type cliCommand struct {
	name		string
	description	string
	callback	func(*config, string) error
}

type Pokemon struct {
	Name		string
	Height		int
	Weight		int
	Stats  map[string]int // Keys will be "hp", "attack", etc.
	Types  []string       // A slice of strings like ["normal", "flying"]
}

type config struct {
	nextLocationsURL *string
	prevLocationsURL *string
	pokeapiClient    *pokeapi.Client
	pokedex			 map[string]Pokemon
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
	"catch": {
		name:		 "catch",
		description: "Catch + pokemon name attempts to catch the pokemon and add them to the user's pokedex",
		callback:	 commandCatch,
	},
	"inspect": {
		name:		 "inspect",
		description: "inspect + caught pokemon name prints out pokemon details",
		callback: 	 commandInspect,
	},
	"pokedex": {
		name:		 "pokedex",
		description: "pokedex prints out a list of caught pokemon",
		callback: 	 commandPokedex,
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
	locationData, err := cfg.pokeapiClient.FetchLocationAreasName(name)
	if err != nil {
		return err
	}
	for _, pokemonEncounters := range locationData.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemonEncounters.Pokemon.Name)
	}
	return nil
}
func commandCatch(cfg *config, name string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	
	pokemonData, err := cfg.pokeapiClient.FetchPokemonData(name)
	if err != nil {
		return err
	}
	
	pokemonExp := pokemonData.BaseExperience
	catchRate := ExpToCatchRate(pokemonExp)
	isCaught := tryCatch(catchRate)
	
	if !isCaught {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}
	
	statsMap := make(map[string]int)
	for _, statData := range pokemonData.Stats {
		statsMap[statData.Stat.Name] = statData.BaseStat
	}

	var typesSlice []string
	for _, typeData := range pokemonData.Types {
		typesSlice = append(typesSlice, typeData.Type.Name)
	}

	caughtPokemon := Pokemon{
		Name:   pokemonData.Name,
		Height: pokemonData.Height,
		Weight: pokemonData.Weight,
		Stats:  statsMap,
		Types:  typesSlice,
	}

	cfg.pokedex[name] = caughtPokemon
	fmt.Printf("%s was caught!\n", name)

	return nil
}
func commandInspect(cfg *config, name string) error {
	// 1. Check if the Pokémon exists in the user's Pokédex
	pokemon, ok := cfg.pokedex[name]
	if !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	// 2. Print the top-level stats
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	// 3. Iterate over the Stats map
	fmt.Println("Stats:")
	orderedStats := []string{"hp", "attack", "defense", "special-attack", "special-defense", "speed"}
	for _, statName := range orderedStats {
    	// Look up the stat directly from the map using the ordered keys
    	if statValue, exists := pokemon.Stats[statName]; exists {
        	fmt.Printf("  -%s: %d\n", statName, statValue)
    	}
	}

	// 4. Iterate over the Types slice
	fmt.Println("Types:")
	for _, typeName := range pokemon.Types {
		fmt.Printf("  - %s\n", typeName)
	}

	return nil
}
func commandPokedex(cfg *config, name string) error {
	fmt.Println("Your Pokedex:")
	
	if len(cfg.pokedex) == 0 {
		fmt.Println(" - You haven't caught any Pokemon yet!")
		return nil
	}

	for _, pokemon := range cfg.pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

// clamp keeps experience stats in the expected range
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ExpToCatchRate converts a Pokemon's base experience into a catch probability.
func ExpToCatchRate(baseExp int) float64 {
	// 1. Define the known boundaries of Pokemon Base Experience
	const minExp = 36.0  // Lowest tier (Magikarp)
	const maxExp = 608.0 // Highest tier (Blissey)

	// 2. Define your desired catch rate boundaries
	// Leaving a small buffer prevents a 100% guarantee or a 0% impossibility.
	const minCatchRate = 0.05  // 5% floor for the hardest encounters
	const maxCatchRate = 0.95  // 95% ceiling for the easiest encounters

	expFloat := float64(baseExp)
	
	// Clamp the input to protect against outliers
	expFloat = clamp(expFloat, minExp, maxExp)

	// 3. Apply the inverse Min-Max Normalization formula
	normalizedExp := (expFloat - minExp) / (maxExp - minExp)
	catchRate := maxCatchRate - (normalizedExp * (maxCatchRate - minCatchRate))

	return catchRate
}

func tryCatch(catchRate float64) bool {
	// rand.Float64() generates a random number between 0.0 and 0.999...
	// If the random number is less than our catch rate, it's a success!
	roll := rand.Float64()
	return roll < catchRate
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
				err := command.callback(cfg, "")
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
