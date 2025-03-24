package main

import "fmt"

func commandScan(config *Config, args []string) error {
	if len(args) != 2 {
		fmt.Println("Expecting: scan <valid path from OneDrive root>")
		return nil
	}
	config.library.ImportFromSource(config.library.sources["onedrive"], args[1])
	return nil
}
