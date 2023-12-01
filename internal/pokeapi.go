package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type GetLocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Location
}
type Args struct {
	Cache *Cache
	Url   string
}

func GetLocation(args Args) GetLocationResponse {
	var url string
	response := GetLocationResponse{}
	if url = args.Url; url == "" {
		url = "https://pokeapi.co/api/v2/location/?offset=0&limit=20"
	}
	if cachedBody, ok := args.Cache.Get(url); ok {
		unerr := json.Unmarshal(cachedBody, &response)
		if unerr != nil {
			panic(unerr)
		}
	} else {
		body := sendRequest(url)
		unerr := json.Unmarshal(body, &response)
		if unerr != nil {
			log.Fatal(unerr)
		}
		args.Cache.Add(url, body)
	}
	return response
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type GetExploreLocationBody struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func GetExploreLocation(context *Context) GetExploreLocationBody {
	response := GetExploreLocationBody{}
	if context.CommandArgs[0] == "" {
		log.Fatal("location is required")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", context.CommandArgs[1])
	fmt.Println(url)

	if cachedBody, ok := context.Cache.Get(url); ok {
		unerr := json.Unmarshal(cachedBody, &response)
		if unerr != nil {
			panic(unerr)
		}
	} else {
		body := sendRequest(url)
		unerr := json.Unmarshal(body, &response)
		if unerr != nil {
			log.Fatal(unerr)
		}
		context.Cache.Add(url, body)
	}
	return response
}

func sendRequest(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status %d", res.StatusCode)
	}
	return body
}
