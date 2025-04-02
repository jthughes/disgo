package player

import (
	"fmt"
	"math/rand/v2"
)

type Album struct {
	Title  string
	Artist string
	Year   string
	Genre  string
	Tracks []Track
}

func (a Album) String() string {
	return fmt.Sprintf("\"%s\" by %s (%d tracks)", a.Title, a.Artist, len(a.Tracks))
}

func (a Album) Play(shuffle bool, repeat bool) {
	fmt.Printf("Album: %s (Shuffle: %t, Repeat: %t)\n", a.String(), shuffle, repeat)

	playlist := make([]int, len(a.Tracks))
	for i := range len(a.Tracks) {
		playlist[i] = i
	}

	if shuffle {
		rand.Shuffle(len(playlist), func(i, j int) {
			// fmt.Printf("[Before] i: %d j: %d", playlist[i], playlist[j])
			playlist[i], playlist[j] = playlist[j], playlist[i]
			// fmt.Printf(" - [After] i: %d j: %d\n", playlist[i], playlist[j])
		})
	}

	for {
		for _, index := range playlist {
			track := a.Tracks[index]
			fmt.Printf("Now playing: %s\n", track.String())
			track.Play()
		}
		if repeat {
			continue
		}
		return
	}

}
