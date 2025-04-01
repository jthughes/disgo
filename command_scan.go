package main

import (
	"fmt"
	"strings"
)

func commandScan(config *Config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Expecting: scan <valid path from OneDrive root>")
		return nil
	}

	path := strings.Join(args, " ")
	config.library.ImportFromSource(config.library.sources["onedrive"], path)
	return nil
}
