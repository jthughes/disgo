package repl

import (
	"fmt"
)

func commandPrevious(config *Config, args []string) error {
	if len(args) != 0 {
		fmt.Println("Expecting: previous")
		return nil
	}
	config.player.Previous()
	return nil
}
