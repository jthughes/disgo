package main

import (
	"fmt"
	"os"
	"strconv"
)

func commandHelp(config *Config, args []string) error {
	fmt.Println("Welcome to GoGroovy!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commands {
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}

// func commandLogin(config *Config, args []string) error {
// 	return nil
// }

func commandScan(config *Config, args []string) error {
	if len(args) != 2 {
		fmt.Println("Expecting: scan <valid path from OneDrive root>")
		return nil
	}
	config.library.ImportFromSource(config.library.sources["onedrive"], args[1])
	return nil
}

func commandList(config *Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Expecting: ls")
		return nil
	}
	albums, err := config.library.GetAlbums()
	if err != nil {
		return fmt.Errorf("failed to get albums: %w", err)
	}
	for i, album := range albums {
		fmt.Printf("[%d] %s (%s): %d tracks\n", i, album.Title, album.Year, len(album.Tracks))
	}
	return nil
}

func commandPlay(config *Config, args []string) error {
	if len(args) != 2 {
		fmt.Println("Expecting: play <#>")
		return nil
	}

	selection, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("unrecognised album number: %w", err)
	}

	albums, err := config.library.GetAlbums()
	if err != nil {
		return fmt.Errorf("failed to get albums: %w", err)
	}

	if selection < 0 || selection >= len(albums) {
		return fmt.Errorf("invalid selection: %d not in range [0, %d]", selection, len(albums)-1)
	}
	album := albums[selection]
	album.Tracks[0].Play()
	return nil
}

func commandExit(config *Config, args []string) error {
	os.Exit(0)
	return nil
}
