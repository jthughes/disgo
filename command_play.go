package main

import (
	"fmt"
	"strconv"
)

func commandPlay(config *Config, args []string) error {
	if len(args) < 2 || len(args) > 3 {
		fmt.Println("Expecting: play <#> [-sr]")
		return nil
	}

	// shuffle := false
	repeat := false

	if len(args) == 3 {
		for i, c := range args[2] {
			if i == 0 {
				if c != '-' {
					return fmt.Errorf("expected '-' got '%c'", c)
				}
				continue
			}
			// if c == 's' {
			// 	shuffle = true
			// } else
			if c == 'r' {
				repeat = true
			} else {
				return fmt.Errorf("unexpected flag '%c'- valid flags: (r)epeat, (s)huffle", c)
			}
		}
	}
	config.player.Repeat = repeat

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
	// album.Play(shuffle, repeat)
	config.player.AddAlbumToPlaylist(album)
	config.player.Play()
	return nil
}
