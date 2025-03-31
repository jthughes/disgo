package main

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

type Player struct {
	Playlist           []Track
	Repeat             bool
	PlaylistPosition   int
	Controller         *beep.Ctrl
	PlaylistCancel     context.CancelFunc
	channelTrackOffset chan int
}

func (p *Player) consumer(ctx context.Context) {
	trackDone := make(chan bool)
	playTrack := func() {
		stream, err := p.Playlist[p.PlaylistPosition].Stream()
		if err != nil {
			fmt.Println("error playing track")
			trackDone <- true
			return
		}
		p.Controller.Streamer = beep.Seq(stream, beep.Callback(func() {
			stream.Close()
			trackDone <- true
		}))
		speaker.Play(p.Controller)
	}
	for {
		select {
		case <-ctx.Done(): // Context cancelled
			return
		case <-trackDone: // Current track is finished
			p.PlaylistPosition += 1
			if p.PlaylistPosition >= len(p.Playlist) {
				if p.Repeat {
					p.PlaylistPosition = 0
				} else {
					p.PlaylistCancel()
					continue
				}
			}
			playTrack()
		case offset := <-p.channelTrackOffset:
			newPosition, err := clamp(p.PlaylistPosition+offset, 0, len(p.Playlist)-1)
			if err != nil {
				return
			}
			speaker.Lock()
			p.Controller.Streamer = nil
			speaker.Unlock()
			p.PlaylistPosition = newPosition
			playTrack()
		}
	}
}

func (p *Player) Init() {
	p.channelTrackOffset = make(chan int)
}

func (p *Player) AddAlbumToPlaylist(a Album) {
	p.Playlist = append(p.Playlist, a.Tracks...)
}

func (p *Player) Play() {
	ctx, cancel := context.WithCancel(context.Background())
	p.PlaylistCancel = cancel
	go p.consumer(ctx)
	p.channelTrackOffset <- 1
}

func (p *Player) Resume() {
	speaker.Lock()
	p.Controller.Paused = false
	speaker.Unlock()
}

func (p *Player) Pause() {
	speaker.Lock()
	p.Controller.Paused = true
	speaker.Unlock()
}

func (p *Player) Next() {
	p.channelTrackOffset <- 1
}

func (p *Player) Previous() {
	p.channelTrackOffset <- -1
}

func (p *Player) Shuffle() {
	rand.Shuffle(len(p.Playlist), func(i, j int) {
		p.Playlist[i], p.Playlist[j] = p.Playlist[j], p.Playlist[i]
	})
}
