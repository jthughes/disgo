package repl

import "fmt"

func commandHelp(config *Config, args []string) error {
	fmt.Println("Welcome to disgo!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commands {
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}
