package main

import (
	"fmt"

	"github.com/gopxl/beep/speaker"
)

func commandStop(config *Config, args []string) error {
	if len(args) != 0 {
		fmt.Println("Expecting: stop")
		return nil
	}
	speaker.Lock()
	config.player.Controller.Streamer = nil
	speaker.Unlock()
	config.player.PlaylistCancel()
	config.player.Playlist = nil
	config.player.PlaylistPosition = -1
	return nil
}
