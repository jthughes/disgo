package main

import (
	"fmt"
)

func commandPrevious(config *Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Expecting: previous")
		return nil
	}
	config.player.Previous()
	return nil
}
