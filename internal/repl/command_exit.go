package repl

import "os"

func commandExit(config *Config, args []string) error {
	os.Exit(0)
	return nil
}
