package main

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

type Player struct {
	Playlist         []Track
	Repeat           bool
	PlaylistPosition int
	Controller       *beep.Ctrl
}

func (p *Player) Init() {
	done := make(chan bool)

	loadStream := func() {
		stream, err := p.Playlist[p.PlaylistPosition].Stream()
		if err != nil {
			fmt.Println("Error")
			return
		}
		p.Controller.Streamer = beep.Seq(stream, beep.Callback(func() {
			stream.Close()
			done <- true
		}))
		speaker.Play(p.Controller)

	}

	go (func() {
		for {
			if len(p.Playlist) > 0 {
				p.JumpPlaylistTo(1)
				loadStream()
				<-done
			}
		}
	})()

}

func (p *Player) AddAlbumToPlaylist(a Album) {
	p.Playlist = append(p.Playlist, a.Tracks...)
}

func (p *Player) AddToPlaylist(t Track) {
	p.Playlist = append(p.Playlist, t)
}

func (p *Player) AddToPlaylistNext(t Track) {
	p.Playlist = slices.Insert(p.Playlist, p.PlaylistPosition+1, t)
}

func (p *Player) PlayNow(t Track) {
	p.Pause()
	p.ClearPlaylist()
	p.AddToPlaylist(t)
	p.Resume()
}

func (p *Player) ClearPlaylist() {
	p.Playlist = nil
	p.StopTrack()
}

func (p *Player) StopTrack() {
	speaker.Lock()
	p.Controller.Streamer = nil
	speaker.Unlock()
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

func (p *Player) PlayNext() {
	p.StopTrack()

}

func (p *Player) Next() {
	p.JumpPlaylistTo(1)
	p.StopTrack()
	p.PlayNext()
}

func (p *Player) Previous() {
	p.JumpPlaylistTo(-1)
}

func (p *Player) JumpPlaylistTo(offset int) {
	p.PlaylistPosition = (p.PlaylistPosition + offset) % len(p.Playlist)
}

func (p *Player) Shuffle() {
	rand.Shuffle(len(p.Playlist), func(i, j int) {
		p.Playlist[i], p.Playlist[j] = p.Playlist[j], p.Playlist[i]
	})
}

type Queue struct {
	streamers []beep.Streamer
}

func (q *Queue) Add(streamers ...beep.Streamer) {
	q.streamers = append(q.streamers, streamers...)
}

func (q *Queue) Stream(samples [][2]float64) (n int, ok bool) {
	// We use the filled variable to track how many samples we've
	// successfully filled already. We loop until all samples are filled.
	filled := 0
	for filled < len(samples) {
		// There are no streamers in the queue, so we stream silence.
		if len(q.streamers) == 0 {
			for i := range samples[filled:] {
				samples[i][0] = 0
				samples[i][1] = 0
			}
			break
		}

		// We stream from the first streamer in the queue.
		n, ok := q.streamers[0].Stream(samples[filled:])
		// If it's drained, we pop it from the queue, thus continuing with
		// the next streamer.
		if !ok {
			q.streamers = q.streamers[1:]
		}
		// We update the number of filled samples.
		filled += n
	}
	return len(samples), true
}

func (q *Queue) Err() error {
	return nil
}
