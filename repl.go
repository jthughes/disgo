package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commands map[string]cliCommand

type History struct {
	commands []string
	position int
}

func (h *History) add(command string) {
	h.commands = append(h.commands, command)
	h.position += 1
}

func (h *History) get() string {
	if len(commands) == 0 {
		return ""
	}
	return h.commands[h.position]
}

func (h *History) prev() string {
	if len(commands) == 0 {
		return ""
	}
	h.position -= 1
	if h.position < 0 {
		h.position = 0
	}
	return h.get()
}

func (h *History) next() string {
	if len(commands) == 0 {
		return ""
	}
	h.position += 1
	if h.position >= len(h.commands) {
		h.position = len(h.commands) - 1
	}
	return h.get()
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

type Config struct {
	library *Library
	player  *Player
}

func registerCommands() (commands map[string]cliCommand) {
	commands = map[string]cliCommand{}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays list of available commands",
		callback:    commandHelp,
	}
	// commands["login"] = cliCommand{
	// 	name:        "login",
	// 	description: "Authenticates with OneDrive",
	// 	callback:    commandLogin,
	// }
	commands["scan"] = cliCommand{
		name:        "scan",
		description: "Scans OneDrive folder for new music",
		callback:    commandScan,
	}
	commands["ls"] = cliCommand{
		name:        "list",
		description: "Lists all albums",
		callback:    commandList,
	}
	commands["play"] = cliCommand{
		name:        "play",
		description: "Plays designated album",
		callback:    commandPlay,
	}
	// commands["next"] = cliCommand{
	// 	name:        "next",
	// 	description: "Plays the next track",
	// 	callback:    commandNext,
	// }
	// commands["previous"] = cliCommand{
	// 	name:        "previous",
	// 	description: "Plays the previous track",
	// 	callback:    commandPrevious,
	// }
	commands["pause"] = cliCommand{
		name:        "pause",
		description: "Pauses the current track",
		callback:    commandPause,
	}
	commands["resume"] = cliCommand{
		name:        "resume",
		description: "Resumes the current track",
		callback:    commandResume,
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	return commands
}

func repl(c *Config) {
	commands = registerCommands()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("GoGroovy > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		command, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := command.callback(c, words)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
