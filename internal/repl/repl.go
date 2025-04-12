package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jthughes/disgo/internal/player"
)

var commands map[string]cliCommand

func processInput(text string) (string, []string) {
	words := strings.Fields(text)
	return words[0], words[1:]
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

type Config struct {
	ConfigPath string
	library    *player.Library
	player     *player.Player
}

func InitConfig(cfgPath string, lib *player.Library, player *player.Player) Config {
	return Config{
		ConfigPath: cfgPath,
		library:    lib,
		player:     player,
	}
}

func registerCommands() (commands map[string]cliCommand) {
	commands = map[string]cliCommand{}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays list of available commands",
		callback:    commandHelp,
	}
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
	commands["playing"] = cliCommand{
		name:        "playing",
		description: "Shows info about current playing track",
		callback:    commandPlaying,
	}
	commands["?"] = cliCommand{
		name:        "?",
		description: "Shows info about current playing track",
		callback:    commandPlaying,
	}
	commands["playlist"] = cliCommand{
		name:        "playlist",
		description: "Shows the current playlist",
		callback:    commandPlaylist,
	}
	commands["stop"] = cliCommand{
		name:        "stop",
		description: "Stops playing current playlist",
		callback:    commandStop,
	}
	commands["repeat"] = cliCommand{
		name:        "repeat",
		description: "Toggles album repeat",
		callback:    commandRepeat,
	}
	commands["next"] = cliCommand{
		name:        "next",
		description: "Plays the next track",
		callback:    commandNext,
	}
	commands["previous"] = cliCommand{
		name:        "previous",
		description: "Plays the previous track",
		callback:    commandPrevious,
	}
	commands["goto"] = cliCommand{
		name:        "goto",
		description: "Plays the track at the given position in the playlist",
		callback:    commandGoto,
	}
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
		description: "Exit disgo",
		callback:    commandExit,
	}
	return commands
}

func Run(c *Config) {
	commands = registerCommands()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("disgo > ")
		scanner.Scan()
		input := scanner.Text()
		commandText, args := processInput(input)
		command, ok := commands[commandText]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := command.callback(c, args)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
