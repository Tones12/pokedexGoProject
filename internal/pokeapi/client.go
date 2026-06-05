package pokeapi

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
)

type pokedexLocArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func FetchLocationAreas(url *string) (pokedexLocArea, error) {
	var locationArea pokedexLocArea
	defaultURL := "https://pokeapi.co/api/v2/location-area/"
	if url == nil {
		url = &defaultURL
	}

	fmt.Println("Making a request to the PokeAPI")
	res, err := http.Get(*url)
	if err != nil {
		return locationArea, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return locationArea, fmt.Errorf("error reading response: %w", err)
	}

	if err := json.Unmarshal(data, &locationArea); err != nil {
		return locationArea, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return locationArea, nil
}