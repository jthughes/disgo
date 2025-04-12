package repl

import (
	"fmt"
)

func commandPlaylist(config *Config, args []string) error {
	if len(args) != 0 {
		fmt.Println("Expecting: playlist")
		return nil
	}
	fmt.Printf("Current Playlist: (Repeat: %t)\n", config.player.Repeat)
	for i, track := range config.player.Playlist {
		out := "   "
		if i == config.player.PlaylistPosition {
			if config.player.Controller.Paused {
				out = " = "
			} else {
				out = " > "
			}
		}
		out += fmt.Sprintf("[%d] ", i)
		out += track.String()
		fmt.Println(out)
	}
	return nil
}
