package repl

import (
	"fmt"
)

func commandResume(config *Config, args []string) error {
	if len(args) != 0 {
		fmt.Println("Expecting: resume")
		return nil
	}
	config.player.Resume()
	return nil
}
