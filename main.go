package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"pokedexcli/internal"
	"strings"
	"time"
)

func commandCatch(context *internal.Context) error {
	areaBody := internal.GetLocationArea(context)
	var caught bool
	source := rand.NewSource(time.Now().UnixNano())
	probability := 15

	for _, result := range areaBody.PokemonEncounters {
		if result.Pokemon.Name == context.CommandArgs[1] {
			context.CatchPokemonUrl = result.Pokemon.Url
			pokemon := internal.GetPokemon(context)
			ranned := rand.New(source)
			roll := ranned.Intn(100)
			if pokemon.BaseExperience < 50 {
				probability = 90
			} else if pokemon.BaseExperience < 100 {
				probability = 70
			} else if pokemon.BaseExperience < 150 {
				probability = 40
			}
			fmt.Println("Throwing a Pokeball at", pokemon.Name, "...")
			if roll < probability {
				caught = true
				fmt.Printf("%s was caught! (%d exp)\n", pokemon.Name, pokemon.BaseExperience)
				context.Pokedex[pokemon.Name] = pokemon
			} else {
				fmt.Println("the pokemon broke free!")
			}
			break
		}
	}
	if !caught {
		fmt.Println("didnt catch anything...")
	}
	return nil
}

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
		"catch": {
			name:        "catch",
			description: "catch a pokemon",
			callback:    commandCatch,
		},
	}

	config := internal.Context{
		Cache:   internal.NewCache(time.Minute * 5),
		Pokedex: make(map[string]internal.Pokemon),
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
