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

func GetLocation(context *Context) GetLocationResponse {
	var url string
	response := GetLocationResponse{}
	if url = context.MapUrl; url == "" {
		url = "https://pokeapi.co/api/v2/location/?offset=0&limit=20"
	}
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
		context.Prev = response.Previous
		context.Next = response.Next
		context.Cache.Add(url, body)
	}
	return response
}

type LocationAreas struct {
	Name string `json:name`
	Url  string `json:url`
}

type GetExploreLocationBody struct {
	Areas []LocationAreas `json:"areas"`
}

func GetExploreLocation(context *Context) GetExploreLocationBody {
	response := GetExploreLocationBody{}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location/%s", context.CommandArgs[1])

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

type FoundPokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon FoundPokemon `json:"pokemon"`
}

type GetLocationAreaBody struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func GetLocationArea(context *Context) GetLocationAreaBody {
	response := GetLocationAreaBody{}
	url := context.LocationAreaUrl
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

type (
	Stat struct {
		Name string `json:"name"`
	}
	Stats struct {
		BaseStat int  `json:"base_stat"`
		Stat     Stat `json:"stat"`
	}
	Type struct {
		Name string `json:"name"`
	}
	Types struct {
		Type Type `json:"type"`
	}
	Pokemon struct {
		Id   int    `json:"id"`
		Name string `json:"name"`

		Height         int     `json:"height"`
		Weight         int     `json:"weight"`
		BaseExperience int     `json:"base_experience"`
		Stats          []Stats `json:"stats"`
		Types          []Types `json:"types"`
	}
)

func GetPokemon(context *Context) Pokemon {
	response := Pokemon{}
	url := context.CatchPokemonUrl

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
