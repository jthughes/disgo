package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/jthughes/disgo/internal/database"
	_ "modernc.org/sqlite"
)

//go:embed sql/schema/001_tracks.sql
var schema string

func main() {

	dbPath := "./library.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Printf("failed to open db: %v", err)
		return
	}

	_, err = db.Exec(schema)
	if err != nil {
		fmt.Printf("Failed to initialize db for library: %v\n", err)
		return
	}

	library := Library{
		dbq:     database.New(db),
		sources: make(map[string]Source),
	}

	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{"User.Read", "Files.Read"},
	}

	source, err := NewOneDriveSource(tokenOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	library.sources[source.String()] = source

	player := Player{
		Playlist:         nil,
		Repeat:           false,
		PlaylistPosition: -1,
		Controller: &beep.Ctrl{
			Streamer: &Queue{},
		},
	}

	config := Config{
		library: &library,
		player:  &player,
	}
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	player.Init()

	repl(&config)
}
