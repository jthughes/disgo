package main

import (
	"fmt"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type AudioMetadata struct {
	Album             string `json:"album"`
	AlbumArtist       string `json:"albumArtist"`
	Artist            string `json:"artist"`
	Bitrate           int    `json:"bitrate"`
	Duration          int    `json:"duration"`
	Genre             string `json:"genre"`
	HasDrm            bool   `json:"hasDrm"`
	IsVariableBitrate bool   `json:"isVariableBitrate"`
	Title             string `json:"title"`
	Track             int    `json:"track"`
	Year              int    `json:"year"`
}

type Track struct {
	Data     File
	FileName string
	Metadata AudioMetadata
	MimeType string
}

func (t Track) Play() error {
	data, err := t.Data.Get()
	if err != nil {
		return err
	}
	var streamer beep.StreamSeekCloser
	var format beep.Format
	switch t.MimeType {
	case "audio/mpeg":
		streamer, format, err = mp3.Decode(data)
	// case "audio/flac":
	// 	streamer, format, err = flac.Decode(data)
	// case "audio/ogg":
	// 	streamer, format, err = vorbis.Decode(data)
	// case "audio/wav":
	// 	streamer, format, err = wav.Decode(data)
	default:
		return fmt.Errorf("unrecognised file type: %s", t.MimeType)
	}

	if err != nil {
		// log.Fatal("couldn't decode file")
		return fmt.Errorf("couldn't decode file: %w", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	return nil
}

func (t Track) Print() {
	fmt.Println(t.String())
}

func (t Track) String() string {
	trackDuration := time.Duration(t.Metadata.Duration * int(time.Millisecond)).Truncate(time.Second).String()
	return fmt.Sprintf("%d - %s (%s)", t.Metadata.Track, t.Metadata.Title, trackDuration)
}
