package pokeapi

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"github.com/tones12/pokedexgoproject/internal/pokecache"
	"time"
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

type pokedexLocAreaName struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Client struct {
    cache *pokecache.Cache
}

func NewClient() *Client {
	timeInterval := 300 * time.Second
	newCache := pokecache.NewCache(timeInterval)
	newClient := Client{
		cache: newCache,
	}
	return &newClient
}

func (c *Client) FetchLocationAreas(url *string) (pokedexLocArea, error) {
	var locationArea pokedexLocArea
	defaultURL := "https://pokeapi.co/api/v2/location-area/"
	if url == nil {
		url = &defaultURL
	}

	cacheData, ok := c.cache.Get(*url)
	if ok {
		if err := json.Unmarshal(cacheData, &locationArea); err != nil {
			return locationArea, fmt.Errorf("error unmarshalling data: %w", err)
		}
		return locationArea, nil
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
	
	c.cache.Add(*url, data)

	if err := json.Unmarshal(data, &locationArea); err != nil {
		return locationArea, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return locationArea, nil
}

func (c *Client) FetchLocationAreasName(name string) (pokedexLocAreaName, error) {
	var locationAreaName pokedexLocAreaName
	defaultURL := "https://pokeapi.co/api/v2/location-area/"
	url := defaultURL + name + "/"

	cacheData, ok := c.cache.Get(url)
	if ok {
		if err := json.Unmarshal(cacheData, &locationAreaName); err != nil {
			return locationAreaName, fmt.Errorf("error unmarshalling data: %w", err)
		}
		return locationAreaName, nil
	}

	fmt.Println("Exploring ")
	res, err := http.Get(*url)
	if err != nil {
		return locationArea, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return locationArea, fmt.Errorf("error reading response: %w", err)
	}
	
	c.cache.Add(*url, data)

	if err := json.Unmarshal(data, &locationArea); err != nil {
		return locationArea, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return locationArea, nil
}