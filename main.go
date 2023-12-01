package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pokedexcli/internal"
	"strings"
	"time"
)

func commandExplore(context *internal.Context) error {
	response := internal.GetExploreLocation(context)
	if context.CommandArgs[0] == "" {
		log.Fatal("location is required")
		return nil
	}
	if len(response.Areas) == 0 {
		fmt.Println("Nothing to explore in", context.CommandArgs[1])
		return nil
	}
	area := response.Areas[0]
	context.LocationAreaUrl = area.Url
	areaBody := internal.GetLocationArea(context)
	fmt.Println("Found Pokemon:")
	for _, results := range areaBody.PokemonEncounters {
		fmt.Println("- ", results.Pokemon.Name)
	}
	return nil
}

func commandMap(config *internal.Context) error {
	config.MapUrl = config.Next
	response := internal.GetLocation(config)
	for _, place := range response.Results {
		fmt.Println(place.Name)
	}

	return nil
}

func commandMapBack(config *internal.Context) error {
	config.MapUrl = config.Prev
	response := internal.GetLocation(config)

	for _, place := range response.Results {
		fmt.Println(place.Name)
	}

	return nil
}

func commandHelp(_ *internal.Context) error {
	fmt.Println("Available commands:")
	fmt.Println("- help: Display this help message")
	fmt.Println("- exit: close the program")
	return nil
}

func commandExit(_ *internal.Context) error {
	fmt.Println("exiting the program")
	os.Exit(0)
	return nil
}

type command struct {
	name        string
	description string
	callback    func(config *internal.Context) error
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	commands := map[string]command{
		"help": {
			name:        "help",
			description: "displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "go forward in the map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "go back in the map",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "explore an area from the map",
			callback:    commandExplore,
		},
	}

	config := internal.Context{
		Cache: internal.NewCache(time.Minute * 5),
	}
	for {
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(input)
		config.CommandArgs = strings.Split(input, " ")
		cmderr := commands[config.CommandArgs[0]].callback(&config)
		if cmderr != nil {
			log.Fatalf("command error %s", cmderr)
		}
	}
}
