package repl

import (
	"github.com/jthughes/disgo/internal/player"
)

type Config struct {
	ConfigPath string
	Library    *player.Library
	Player     *player.Player
}

func InitConfig(cfgPath string, lib *player.Library, player *player.Player) Config {
	return Config{
		ConfigPath: cfgPath,
		Library:    lib,
		Player:     player,
	}
}
