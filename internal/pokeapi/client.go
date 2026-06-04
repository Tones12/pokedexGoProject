package pokeapi

import (
	"fmt"
	"net/http"
)

type config struct {

}

func FetchLocationAreas() error {
	fmt.Println("Making a request to the PokeAPI")
		res, err := http.Get("https://pokeapi.co/api/v2/location-area/")
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if err := json.Unmarshal(data, &variable); err != nil {
		return err
	}
}