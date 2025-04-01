package main

import "fmt"

func commandList(config *Config, args []string) error {
	if len(args) != 0 {
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
