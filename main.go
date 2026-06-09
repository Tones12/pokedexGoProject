package main

import "github.com/tones12/pokedexgoproject/internal/pokeapi"

func main() {
	cfg := config{
		pokeapiClient: pokeapi.NewClient(),
	}
	startRepl(&cfg)
}