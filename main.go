package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pokedexcli/pokeapi"
	"strings"
)

func commandMap(config map[string]string) error {
	args := map[string]string{"url": config["next"]}
	response := pokeapi.GetLocation(args)
	for _, place := range response.Results {
		fmt.Println(place.Name)
	}
	config["prev"] = response.Previous
	config["next"] = response.Next

	return nil
}

func commandMapBack(config map[string]string) error {
	args := map[string]string{"url": config["prev"]}
	response := pokeapi.GetLocation(args)
	for _, place := range response.Results {
		fmt.Println(place.Name)
	}
	config["prev"] = response.Previous
	config["next"] = response.Next
	return nil
}

func commandHelp(_ map[string]string) error {
	fmt.Println("Available commands:")
	fmt.Println("- help: Display this help message")
	fmt.Println("- exit: close the program")
	return nil
}

func commandExit(_ map[string]string) error {
	fmt.Println("exiting the program")
	os.Exit(0)
	return nil
}

type command struct {
	name        string
	description string
	callback    func(config map[string]string) error
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

	config := make(map[string]string)
	for {
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(input)
		cmderr := commands[input].callback(config)
		if cmderr != nil {
			log.Fatalf("command error %s", cmderr)
		}
	}
}
