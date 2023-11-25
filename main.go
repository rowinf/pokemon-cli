package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func commandHelp() error {
	fmt.Println("Available commands:")
	fmt.Println("- help: Display this help message")
	fmt.Println("- exit: close the program")
	return nil
}

func commandExit() error {
	fmt.Println("exiting the program")
	os.Exit(0)
	return nil
}

type command struct {
	name        string
	description string
	callback    func() error
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
	}

	for {
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(input)
		commands[input].callback()
	}
}
