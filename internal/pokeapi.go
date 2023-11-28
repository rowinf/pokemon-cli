package internal

import (
	"encoding/json"
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
		unerr := json.Unmarshal(body, &response)
		if unerr != nil {
			log.Fatal(unerr)
		}
		args.Cache.Add(url, body)
	}
	return response
}
