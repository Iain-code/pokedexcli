package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(client *pokeapi.Client, area string) error
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},

	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},

	"map": {
		name:        "map",
		description: "Displays the map",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays previous map",
		callback:    commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "Reveals pokemon in the area",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "gotta catch em all",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "check out that pokemon",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "show all caught pokemon",
		callback:    commandPokedex,
	},
}

func main() {

	cache := pokecache.NewCache(5 * time.Second) // create once

	client := pokeapi.NewClient(cache) // pass cache to client

	scanner := bufio.NewScanner(os.Stdin)
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("Welcome to the Pokedex!\n")
	client.CaughtPokemon = make(map[string]*pokeapi.Pokemon)
	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		cmd := ""
		input := scanner.Text()
		input = strings.ToLower(input)
		if len(input) == 0 {
			fmt.Println("Please input command")
		} else {
			trimmed := strings.Fields(input)

			if trimmed[0] == "explore" || trimmed[0] == "catch" || trimmed[0] == "inspect" {
				if len(trimmed) > 1 {
					cmd = trimmed[1]
				} else {
					fmt.Println("Please enter command")
				}
			}

			switch trimmed[0] {
			case "exit":
				commands["exit"].callback(client, cmd)
			case "help":
				commands["help"].callback(client, cmd)
			case "map":
				commands["map"].callback(client, cmd)
			case "mapb":
				commands["mapb"].callback(client, cmd)
			case "explore":
				commands["explore"].callback(client, cmd)
			case "catch":
				commands["catch"].callback(client, cmd)
			case "inspect":
				commands["inspect"].callback(client, cmd)
			case "pokedex":
				commands["pokedex"].callback(client, cmd)
			default:
				continue
			}
		}
	}
}

func commandExit(client *pokeapi.Client, area string) error {

	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(client *pokeapi.Client, area string) error {

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandMap(client *pokeapi.Client, area string) error {

	url := "https://pokeapi.co/api/v2/location-area/"

	if client.Next != nil {
		url = *client.Next
	}

	locationrequest, err := client.GetLocationAreas(url)
	if err != nil {
		return errors.New("location area not recieved")
	}

	for _, location := range locationrequest.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(client *pokeapi.Client, area string) error {

	if client.Previous != nil {

	}
	locationrequest, err := client.GetLocationAreas(*client.Previous)
	if err != nil {
		return nil
	}

	for _, location := range locationrequest.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(client *pokeapi.Client, inputarea string) error {

	if len(inputarea) <= 0 {
		fmt.Println("No area given")
	}

	url := "https://pokeapi.co/api/v2/location-area/"

	areaUrl := inputarea

	datarequest, err := client.GetLocationAreas(url)

	if err != nil {
		return fmt.Errorf("error is not nil")
	}

	for _, area := range datarequest.Results {
		if area.Name == inputarea {
			areaUrl = area.URL
		}
	}

	newrequest, err := client.GetLocationAreaName(areaUrl)
	if err != nil {
		return fmt.Errorf("error is not nil")
	}
	fmt.Printf("Exploring %s...\n", inputarea)
	for _, data := range newrequest.PokemonEncounters {

		fmt.Println(data.Pokemon.Name)
	}
	return nil
}

func commandCatch(client *pokeapi.Client, pokemon string) error {

	if pokemon == "" || pokemon == " " {
		fmt.Println("Please enter pokemon name")
	} else {
		fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
		pokemonUrl := "https://pokeapi.co/api/v2/pokemon/" + pokemon

		pokemonData, err := client.GetPokemonData(pokemonUrl)
		if err != nil {
			return err
		}
		randomNum := rand.Intn(350)

		if randomNum > pokemonData.BaseExperience {
			fmt.Printf("%s was caught\n", pokemon)
			fmt.Printf("You may now inspect it with the inspect command.\n")
			_, exists := client.CaughtPokemon[pokemon]
			if !exists {
				client.CaughtPokemon[pokemon] = pokemonData
			} else {
				fmt.Println("We already have that Pokemon, no point keeping another")
			}
		} else {
			fmt.Printf("You throw like Arron and the %s escaped\n", pokemon)
		}
	}
	return nil
}

func commandInspect(client *pokeapi.Client, pokemon string) error {

	// use the map we made and check if pokemon name is in it
	// return info if it is or say no

	value, exists := client.CaughtPokemon[pokemon]
	if !exists {
		fmt.Printf("You havnt caught a %s yet...go and catch one first!!\n", pokemon)
	} else {
		fmt.Printf("Name: %s\n", value.Name)
		fmt.Printf("Experience: %d\n", value.BaseExperience)

		for i, ability := range value.Abilities {
			fmt.Printf("Ability %d: %v\n", i+1, ability.Ability.Name)
		}
		fmt.Printf("")
	}
	return nil
}

func commandPokedex(client *pokeapi.Client, pokemon string) error {

	for _, pokemon := range client.CaughtPokemon {
		fmt.Println("-", pokemon.Name)
	}
	return nil
}
