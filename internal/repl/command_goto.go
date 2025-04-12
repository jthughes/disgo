package repl

import (
	"fmt"
	"strconv"
)

func commandGoto(config *Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Expecting: goto <#>")
		return nil
	}
	targetPosition, err := strconv.Atoi(args[0])
	if err != nil || targetPosition < 0 || targetPosition > len(config.player.Playlist) {
		fmt.Printf("Error: Needs to be in range [0, %d]", len(config.player.Playlist)-1)
		return nil
	}
	offset := targetPosition - config.player.PlaylistPosition
	config.player.JumpTo(offset)
	return nil
}
