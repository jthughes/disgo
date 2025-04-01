package main

import (
	"fmt"
)

func commandPlaying(config *Config, args []string) error {
	if len(args) != 0 {
		fmt.Println("Expecting: 'playing' or '?'")
		return nil
	}
	if config.player.Controller.Streamer != nil && len(config.player.Playlist) > 0 && config.player.PlaylistPosition >= 0 {
		out := fmt.Sprintf("Now Playing: %s", config.player.Playlist[config.player.PlaylistPosition].String())
		if config.player.Controller.Paused {
			out += " - Paused"
		}
		fmt.Println(out)
	} else {
		fmt.Printf("Now Playing: Nothing\n")
	}
	return nil
}
