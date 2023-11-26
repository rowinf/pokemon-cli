package pokeapi

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
type ApiResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Location
}

func GetLocation(args map[string]string) ApiResponse {
	url := "https://pokeapi.co/api/v2/location/"
	if argUrl, ok := args["url"]; ok && argUrl != "" {
		url = argUrl
	}
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
	// fmt.Printf("%s", body)
	response := ApiResponse{}
	unerr := json.Unmarshal(body, &response)
	if unerr != nil {
		log.Fatal(unerr)
	}
	return response
}
