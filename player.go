package main

import (
	"io"
	"log"
	"time"

	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func play(data io.ReadCloser) {
	streamer, format, err := mp3.Decode(data)
	if err != nil {
		log.Fatal("couldn't decode file")
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
	select {}
}
