package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	pokeapi "pokedexcli/internal"
	"strings"
	"time"
)

func commandMap(config *Context) error {
	args := pokeapi.Args{}
	args.Url = config.Next
	args.Cache = config.Cache
	response := pokeapi.GetLocation(args)
	for _, place := range response.Results {
		fmt.Println(place.Name)
	}
	config.Prev = response.Previous
	config.Next = response.Next

	return nil
}

func commandMapBack(config *Context) error {
	args := pokeapi.Args{}
	args.Url = config.Prev
	args.Cache = config.Cache
	response := pokeapi.GetLocation(args)

	for _, place := range response.Results {
		fmt.Println(place.Name)
	}
	config.Prev = response.Previous
	config.Next = response.Next
	return nil
}

func commandHelp(_ *Context) error {
	fmt.Println("Available commands:")
	fmt.Println("- help: Display this help message")
	fmt.Println("- exit: close the program")
	return nil
}

func commandExit(_ *Context) error {
	fmt.Println("exiting the program")
	os.Exit(0)
	return nil
}

type Context struct {
	Prev  string
	Next  string
	Cache *pokeapi.Cache
}

type command struct {
	name        string
	description string
	callback    func(config *Context) error
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
	}

	config := Context{}
	cache := pokeapi.NewCache(time.Minute * 5)
	config.Cache = cache
	for {
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(input)
		cmderr := commands[input].callback(&config)
		if cmderr != nil {
			log.Fatalf("command error %s", cmderr)
		}
	}
}
