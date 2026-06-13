package main

import "github.com/tones12/pokedexgoproject/internal/pokeapi"

func main() {
	cfg := &config{
		pokeapiClient: pokeapi.NewClient(),
		pokedex: make(map[string]Pokemon),
	}
	startRepl(cfg)
}