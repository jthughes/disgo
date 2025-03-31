package main

import (
	"fmt"
)

func commandRepeat(config *Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Expecting: repeat")
		return nil
	}
	config.player.Repeat = !config.player.Repeat
	fmt.Printf("Repeat set to: %t\n", config.player.Repeat)
	return nil
}
